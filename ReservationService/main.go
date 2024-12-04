// package main

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"

// 	_ "github.com/go-sql-driver/mysql"
// )

// // Global db connection
// var db *sql.DB

// // VehicleSlot struct holds the details for each vehicle slot
// type VehicleSlot struct {
// 	VehicleID              string
// 	AvailableSlotStartDate string
// 	AvailableSlotEndDate   string
// 	AvailableSlotStartTime string
// 	AvailableSlotEndTime   string
// }

// // Initialize the database connection
// func initDB() {
// 	dsn := "user:password@tcp(127.0.0.1:3306)/CarSharing" // Replace with your MySQL credentials
// 	var err error
// 	db, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Ensure the database connection is valid
// 	if err = db.Ping(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func main() {
// 	// Initialize the database
// 	initDB()
// 	defer db.Close()

// 	// Set up routes
// 	http.HandleFunc("/reserve", reserveSlot) // This handles slot reservation
// 	http.HandleFunc("/", showAvailableSlots) // This shows available vehicle slots

//		// Start the server
//		log.Fatal(http.ListenAndServe(":8082", nil))
//	}
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

	// Set up routes
	http.HandleFunc("/", showAvailableSlots)                     // Displays available slots
	http.HandleFunc("/reserve", reserveSlot)                     // Reserves a slot and redirects
	http.HandleFunc("/modify-booking", showModifyBookingHandler) // Modify booking form
	//http.HandleFunc("/update-booking", modifyBookingHandler)     // Update the booking

	// Start the server
	log.Println("Server is running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
