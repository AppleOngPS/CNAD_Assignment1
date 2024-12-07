package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Global db connection
var db *sql.DB

// User struct to hold user data (Exported to be accessible in other files)
type User struct {
	UserID                string
	Username              string
	Email                 string
	Password              string
	MembershipID          string
	MembershipStatus      string
	MembershipDescription string
	MembershipDiscount    string
}

// Initialize the database connection
func initDB() {
	dsn := "user:password@tcp(127.0.0.1:3306)/CarSharing" // Replace with your MySQL credentials
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the database connection is valid
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	// Set up routes
	// Signup and login pages
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/signup/submit", signup)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/login/submit", login)

	// Profile page
	http.HandleFunc("/profile", profilePage)
	http.HandleFunc("/profile/update", updateProfile)
	http.HandleFunc("/profile/delete", deleteProfile)
	http.HandleFunc("/verify", verifyEmail)
	// Car listing page
	//http.HandleFunc("/car-listing", func(w http.ResponseWriter, r *http.Request) {
	//CarListingHandler(w, r, db) // Call to carListing.go handler function
	//})

	// Car Reservation page
	//http.HandleFunc("/reservation", Reservation)
	//http.HandleFunc("/reservation/submit", reservationHandler)
	//http.HandleFunc("/reservation/success", reservationSuccess)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
