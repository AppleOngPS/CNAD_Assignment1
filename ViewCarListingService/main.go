package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Global db connection
var db *sql.DB

type Car struct {
	VehicleID    string
	VehicleBrand string
	StartDate    string
	EndDate      string
	StartTime    string
	EndTime      string
	Amount       float64
	Location     string
	ChargeLevel  string
	Cleanliness  string
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

	// Car listing page
	http.HandleFunc("/car-listing", func(w http.ResponseWriter, r *http.Request) {
		CarListingHandler(w, r, db) // Call to carListing.go handler function
	})

	// Start the server
	log.Fatal(http.ListenAndServe(":8081", nil))
}
