package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MakMoinee/go-mith/pkg/email"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Handle the signup page
func signupPage(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(tmpl))
}

// Generate a unique verification token
func generateVerificationToken() string {
	return uuid.New().String() // Generate a unique token using UUID
}

// Send the verification email
func sendVerificationEmail(userEmail, token string) error {
	// Construct the verification URL
	verificationURL := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)

	// Construct the email content
	subject := "Email Verification"
	body := fmt.Sprintf("Please click the following link to verify your email: %s", verificationURL)

	// Use your email service to send the email
	emailService := email.NewEmailService(587, "smtp.gmail.com", "ongapple1@gmail.com", "dsvw cwzk oieg nglb")
	_, err := emailService.SendEmail(userEmail, subject, body)
	return err
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

		// Generate a unique verification token
		token := generateVerificationToken()

		// Insert the user into the database with 'is_verified' flag set to FALSE
		_, err = db.Exec(`
			INSERT INTO users (username, email, password, membershipID, is_verified, verification_token) 
			VALUES (?, ?, ?, ?, FALSE, ?)`,
			username, email, string(hashedPassword), membershipID, token,
		)

		if err != nil {
			http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Send verification email
		err = sendVerificationEmail(email, token)
		if err != nil {
			http.Error(w, "Failed to send verification email: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to the login page after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Handle email verification
func verifyEmail(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the URL query
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}

	// Retrieve the user by token
	var userID int
	var isVerified bool
	err := db.QueryRow(`
		SELECT userid, is_verified
		FROM users
		WHERE verification_token = ?`, token).Scan(&userID, &isVerified)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid verification token", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the email is already verified
	if isVerified {
		http.Error(w, "Email is already verified", http.StatusBadRequest)
		return
	}

	// Update the user's status to verified
	_, err = db.Exec(`
		UPDATE users
		SET is_verified = TRUE
		WHERE userid = ?`, userID)
	if err != nil {
		http.Error(w, "Failed to update verification status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Success: Redirect to the login page or confirmation page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
