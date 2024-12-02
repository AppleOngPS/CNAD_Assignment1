package main

import (
	"fmt"
	"net/http"
)

// Serve the Profile page with Update and Delete options
func profilePage(w http.ResponseWriter, r *http.Request) {
	// Assume user ID is passed as a query parameter
	userID := r.URL.Query().Get("userID")

	// Query the database to get the user details
	var user User
	err := db.QueryRow("SELECT userID, username, email, password, membershipID FROM users WHERE userID = ?", userID).
		Scan(&user.UserID, &user.Username, &user.Email, &user.Password, &user.MembershipID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// HTML for Profile Page with Update and Delete options
	tmpl := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Profile</title>
		</head>
		<body>
			<h2>Your Profile</h2>
			<p><strong>Username:</strong> %s</p>
			<p><strong>Email:</strong> %s</p>
			<p><strong>Membership ID:</strong> %s</p>

			<h3>Update your details:</h3>
			<form action="/profile/update" method="POST">
				<input type="hidden" name="userID" value="%s">
				<label for="username">Username:</label>
				<input type="text" id="username" name="username" value="%s" required><br><br>

				<label for="email">Email:</label>
				<input type="email" id="email" name="email" value="%s" required><br><br>

				<label for="password">Password:</label>
				<input type="password" id="password" name="password" value="%s" required><br><br>

				<input type="submit" value="Update Profile">
			</form>

			<h3>Delete your account:</h3>
			<form action="/profile/delete" method="POST">
				<input type="hidden" name="userID" value="%s">
				<input type="submit" value="Delete Account" onclick="return confirm('Are you sure you want to delete your account?');">
			</form>
		</body>
		</html>
	`, user.Username, user.Email, user.MembershipID, user.UserID, user.Username, user.Email, user.Password, user.UserID)

	// Serve the HTML to the user
	w.Write([]byte(tmpl))
}

// Handle the update profile submission
func updateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get the user input from the form
		userID := r.FormValue("userID")
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Update the user's details in the database
		_, err := db.Exec("UPDATE users SET username = ?, email = ?, password = ? WHERE userID = ?", username, email, password, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the profile page after successful update
		http.Redirect(w, r, "/profile?userID="+userID, http.StatusSeeOther)
	}
}

// Handle the delete account action
// Handle the delete action
func deleteProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get userID from the form submission
		userID := r.FormValue("userID")

		// Check if userID is empty or not
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		// Perform the DELETE operation in the database
		_, err := db.Exec("DELETE FROM users WHERE userID = ?", userID)
		if err != nil {
			http.Error(w, "Failed to delete the user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the login page after successful deletion
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
