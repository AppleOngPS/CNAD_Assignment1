package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// User struct to hold the data from the database
type User struct {
	UserID       string
	Username     string
	Email        string
	Password     string
	MembershipID string
}

// Retrieve data from the database
func GetData(db *sql.DB) {
	results, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	fmt.Println("Data in the table:")
	for results.Next() {
		var p User
		err = results.Scan(&p.UserID, &p.Username, &p.Email, &p.Password, &p.MembershipID)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(p.UserID, p.Username, p.Email, p.Password, p.MembershipID)
	}
}

// Insert data into the database
func InsertData(db *sql.DB) {
	p := User{UserID: "1", Username: "John", Email: "john@gmail.com", Password: "1234", MembershipID: "VIP"}
	result, err := db.Exec("INSERT INTO users (userID, username, email, password, membershipID) VALUES (?, ?, ?, ?, ?)", p.UserID, p.Username, p.Email, p.Password, p.MembershipID)
	if err != nil {
		fmt.Printf("Error inserting data: %v\n", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error fetching affected rows: %v\n", err)
		return
	}
	fmt.Printf("Rows affected: %d\n", rowsAffected)
}

// Delete data from the database based on the ID
func DeleteData(db *sql.DB, userID string) {
	// First, delete the records in the billing table that reference the reservationID
	_, err := db.Exec("DELETE FROM billing WHERE reservationID IN (SELECT reservationID FROM reservation WHERE userID = ?)", userID)
	if err != nil {
		fmt.Printf("Error deleting data from billing table: %v\n", err)
		return
	}

	// delete the records in the reservation table that reference to the user table
	_, err = db.Exec("DELETE FROM reservation WHERE userID = ?", userID)
	if err != nil {
		fmt.Printf("Error deleting data from reservation table: %v\n", err)
		return
	}

	// delete the records in the trackrentalhistory table that reference to the user table
	_, err = db.Exec("DELETE FROM trackrentalhistory WHERE userID = ?", userID)
	if err != nil {
		fmt.Printf("Error deleting data from trackrentalhistory table: %v\n", err)
		return
	}

	// delete the user record from the users table
	result, err := db.Exec("DELETE FROM users WHERE userID = ?", userID)
	if err != nil {
		fmt.Printf("Error deleting data from users table: %v\n", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error fetching affected rows: %v\n", err)
		return
	}

	if rowsAffected > 0 {
		fmt.Printf("Successfully deleted record with ID: %s\n", userID)
	} else {
		fmt.Printf("No record found with ID: %s\n", userID)
	}
}

func main() {
	// Update with your database credentials
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/CarSharing")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Ensure the database connection is valid
	if err = db.Ping(); err != nil {
		panic(err.Error())
	}

	// Insert data
	InsertData(db)

	// Fetch data to see the record in the table
	GetData(db)

	// Delete data based on ID
	DeleteData(db, "U1") // Replace "1" with the ID of the record you want to delete

	// Fetch data again to confirm deletion
	GetData(db)
}
