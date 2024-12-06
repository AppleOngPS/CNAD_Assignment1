package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Handle the login page
func loginPage(w http.ResponseWriter, r *http.Request) {
	// Embed the login HTML form
	tmpl := `<!DOCTYPE html>
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
    </html>`

	// Serve the HTML to the user
	w.Write([]byte(tmpl))
}

// Handle the login form submission
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Retrieve user from the database based on email
		var user User
		var storedHashedPassword string
		err := db.QueryRow("SELECT userID, username, email, password, membershipID FROM users WHERE email = ?", email).Scan(&user.UserID, &user.Username, &user.Email, &storedHashedPassword, &user.MembershipID)

		// If no user is found with the provided email
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Compare the entered password with the stored hashed password
		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
		if err != nil {
			// Password does not match
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Redirect to the profile page after successful login
		http.Redirect(w, r, "/profile?userID="+user.UserID, http.StatusSeeOther)
	}
}
