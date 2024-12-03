package main

import (
	"net/http"
)

// Handle the signup page
func signupPage(w http.ResponseWriter, r *http.Request) {
	// HTML for the signup form with a membership dropdown, defaulting to "Basic"
	tmpl := `
		<!DOCTYPE html>
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
				<input type="hidden" name="membershipID" value="1"><br><br>

				<input type="submit" value="Sign Up">
			</form>
		</body>
		</html>
	`

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

		// Fetch the membership description for the selected membershipID (Basic is always ID=1)
		var membershipDescription string
		err := db.QueryRow("SELECT descriptions FROM membership WHERE membershipID = ?", membershipID).Scan(&membershipDescription)
		if err != nil {
			http.Error(w, "Failed to retrieve membership description: "+err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
	INSERT INTO users (username, email, password, typeOfStatus) 
	VALUES (?, ?, ?, (SELECT typeOfStatus FROM membership WHERE membershipID = ?))`,
			username, email, password, membershipID,
		)

		if err != nil {
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the login page after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
