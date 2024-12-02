// package main

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/go-sql-driver/mysql"
// )

// // User struct to hold the data from the database
// type User struct {
// 	UserID       string
// 	Username     string
// 	Email        string
// 	Password     string
// 	MembershipID string
// }

// // Retrieve data from the database
// func GetData(db *sql.DB) {
// 	results, err := db.Query("SELECT * FROM users")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer results.Close()

// 	fmt.Println("Data in the table:")
// 	for results.Next() {
// 		var p User
// 		err = results.Scan(&p.UserID, &p.Username, &p.Email, &p.Password, &p.MembershipID)
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		fmt.Println(p.UserID, p.Username, p.Email, p.Password, p.MembershipID)
// 	}
// }

// // Insert data into the database
// func InsertData(db *sql.DB) {
// 	p := User{UserID: "U5", Username: "John", Email: "john@gmail.com", Password: "1234", MembershipID: "M1"}
// 	result, err := db.Exec("INSERT INTO users (userID, username, email, password, membershipID) VALUES (?, ?, ?, ?, ?)", p.UserID, p.Username, p.Email, p.Password, p.MembershipID)
// 	if err != nil {
// 		fmt.Printf("Error inserting data: %v\n", err)
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		fmt.Printf("Error fetching affected rows: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("Rows affected: %d\n", rowsAffected)
// }

// // Update data in the users table
// func UpdateData(db *sql.DB, userID, username, email, password, membershipID string) {
// 	// Update user details in the users table
// 	result, err := db.Exec("UPDATE users SET username = ?, email = ?, password = ?, membershipID = ? WHERE userID = ?", username, email, password, membershipID, userID)
// 	if err != nil {
// 		fmt.Printf("Error updating data: %v\n", err)
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		fmt.Printf("Error fetching affected rows: %v\n", err)
// 		return
// 	}

// 	if rowsAffected > 0 {
// 		fmt.Printf("Successfully updated user with ID: %s\n", userID)
// 	} else {
// 		fmt.Printf("No record found with ID: %s\n", userID)
// 	}
// }

// // Delete data from the database based on the ID
// func DeleteData(db *sql.DB, userID string) {
// 	// First, delete the records in the billing table that reference the reservationID
// 	_, err := db.Exec("DELETE FROM billing WHERE reservationID IN (SELECT reservationID FROM reservation WHERE userID = ?)", userID)
// 	if err != nil {
// 		fmt.Printf("Error deleting data from billing table: %v\n", err)
// 		return
// 	}

// 	// Then, delete the records in the reservation table that reference the userID
// 	_, err = db.Exec("DELETE FROM reservation WHERE userID = ?", userID)
// 	if err != nil {
// 		fmt.Printf("Error deleting data from reservation table: %v\n", err)
// 		return
// 	}

// 	// Now, delete the records in the trackrentalhistory table that reference the userID
// 	_, err = db.Exec("DELETE FROM trackrentalhistory WHERE userID = ?", userID)
// 	if err != nil {
// 		fmt.Printf("Error deleting data from trackrentalhistory table: %v\n", err)
// 		return
// 	}

// 	// Finally, delete the user record from the users table
// 	result, err := db.Exec("DELETE FROM users WHERE userID = ?", userID)
// 	if err != nil {
// 		fmt.Printf("Error deleting data from users table: %v\n", err)
// 		return
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		fmt.Printf("Error fetching affected rows: %v\n", err)
// 		return
// 	}

// 	if rowsAffected > 0 {
// 		fmt.Printf("Successfully deleted record with ID: %s\n", userID)
// 	} else {
// 		fmt.Printf("No record found with ID: %s\n", userID)
// 	}
// }

// func main() {
// 	// Update with your database credentials
// 	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/CarSharing")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer db.Close()

// 	// Ensure the database connection is valid
// 	if err = db.Ping(); err != nil {
// 		panic(err.Error())
// 	}

// 	// Insert data
// 	InsertData(db)

// 	// Fetch data to see the record in the table
// 	GetData(db)

// 	// Update data based on ID
// 	UpdateData(db, "U3", "Mary", "m@gmail.com", "new123", "m2")

// 	// Fetch data again to confirm update
// 	GetData(db)

// 	// Delete data based on ID
// 	DeleteData(db, "U1")

// 	// Fetch data again to confirm deletion
// 	GetData(db)
// }

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
	UserID       string
	Username     string
	Email        string
	Password     string
	MembershipID string
}

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
	// Signup and login pages
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/signup/submit", signup)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/login/submit", login)

	// Profile page
	http.HandleFunc("/profile", profilePage)
	http.HandleFunc("/profile/update", updateProfile)
	http.HandleFunc("/profile/delete", deleteProfile)

	// Car listing page
	http.HandleFunc("/car-listing", func(w http.ResponseWriter, r *http.Request) {
		CarListingHandler(w, r, db) // Call to carListing.go handler function
	})

	// Car Reservation page
	http.HandleFunc("/reservation", Reservation)
	//http.HandleFunc("/reservation/submit", reservationHandler)
	//http.HandleFunc("/reservation/success", reservationSuccess)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
