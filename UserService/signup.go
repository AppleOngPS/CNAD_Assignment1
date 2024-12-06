package main

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Handle the signup page
func signupPage(w http.ResponseWriter, r *http.Request) {
	// HTML for the signup form with a membership dropdown, defaulting to "Basic"
	tmpl := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Sign Up</title>
	</head>
	<body>
		<h2>Create Account</h2>
		<form action="/signup/submit" method="POST">
			<label for="username">Username:</label>
			<input type="text" id="username" name="username" required><br><br>

			<label for="email">Email:</label>
			<input type="email" id="email" name="email" required><br><br>

			<label for="password">Password:</label>
			<input type="password" id="password" name="password" required><br><br>

			<!-- Default "Basic" membership selection -->
			<input type="hidden" name="membershipID" value="M1"><br><br>

			<input type="submit" value="Sign Up">
		</form>
	</body>
	</html>`

	// Serve the HTML to the user
	w.Write([]byte(tmpl))
}

// Handle the signup form submission
func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		membershipID := r.FormValue("membershipID")

		// Hash the password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch the membership description for the selected membershipID
		var membershipDescription string
		err = db.QueryRow("SELECT memberDescriptions FROM membership WHERE membershipID = ?", membershipID).Scan(&membershipDescription)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid membership ID", http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, "Failed to retrieve membership description: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Insert the user into the database
		_, err = db.Exec(`
			INSERT INTO users (username, email, password, membershipID) 
			VALUES (?, ?, ?, ?)`,
			username, email, string(hashedPassword), membershipID,
		)

		if err != nil {
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the login page after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
