package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Global db connection
var db *sql.DB

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
	// to run this page put url like like this http://localhost:8083/displayCurrentbooking?reservationID=10 in browser reservationid can be changed

	// Set up routes
	http.HandleFunc("/confirmReservation", confirmReservation)
	http.HandleFunc("/currentBooking", displayCurrentBooking)

	// Start the server
	log.Println("Server is running on port 8083...")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
