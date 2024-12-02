package main

import (
	"fmt"
	"net/http"
)

// Serve the profile page after login (Display user data)
func profilePage(w http.ResponseWriter, r *http.Request) {
	// Get the userID from the session cookie
	cookie, err := r.Cookie("userID")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Query the database for user details using userID from cookie
	var user User
	err = db.QueryRow("SELECT username, email, membershipID FROM users WHERE userID = ?", cookie.Value).
		Scan(&user.Username, &user.Email, &user.MembershipID)
	if err != nil {
		http.Error(w, "Could not fetch user details", http.StatusInternalServerError)
		return
	}

	// HTML page embedded directly in the Go code
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>Your Profile</title>
		</head>
		<body>
		    <h2>Your Profile</h2>
		    <p><strong>Username:</strong> %s</p>
		    <p><strong>Email:</strong> %s</p>
		    <p><strong>Membership Type:</strong> %s</p>
		    <a href="/logout">Logout</a>
		</body>
		</html>`, user.Username, user.Email, user.MembershipID)

	// Write the HTML to the response
	w.Write([]byte(html))
}

// Logout the user
func logout(w http.ResponseWriter, r *http.Request) {
	// Remove the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "userID",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
