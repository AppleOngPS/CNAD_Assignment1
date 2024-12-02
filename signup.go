package main

import (
	"net/http"
)

// Handle the signup page
func signupPage(w http.ResponseWriter, r *http.Request) {
	// Embed the signup HTML form
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

				<input type="hidden" name="membership" value="M1"><br><br>

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
		membershipID := r.FormValue("membership")

		// Insert the new user into the database
		_, err := db.Exec("INSERT INTO users (username, email, password, membershipID) VALUES (?, ?, ?, ?)",
			username, email, password, membershipID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Redirect to the login page after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
