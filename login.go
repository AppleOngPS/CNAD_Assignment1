package main

import (
	"net/http"
)

// Handle the login page
func loginPage(w http.ResponseWriter, r *http.Request) {
	// Embed the login HTML form
	tmpl := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Login</title>
		</head>
		<body>
			<h2>Login</h2>
			<form action="/login/submit" method="POST">
				<label for="email">Email:</label>
				<input type="email" id="email" name="email" required><br><br>

				<label for="password">Password:</label>
				<input type="password" id="password" name="password" required><br><br>

				<input type="submit" value="Login">
			</form>
		</body>
		</html>
	`
	// Serve the HTML to the user
	w.Write([]byte(tmpl))
}

// Handle the login form submission
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Retrieve user from the database based on email and password
		var user User
		err := db.QueryRow("SELECT userID, username, email, password, membershipID FROM users WHERE email = ? AND password = ?", email, password).Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.MembershipID)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set a session cookie for the user
		http.SetCookie(w, &http.Cookie{
			Name:  "userID",
			Value: user.UserID,
			Path:  "/",
		})

		// Redirect to the profile page after successful login
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}
