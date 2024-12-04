// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"strings"

// 	_ "github.com/go-sql-driver/mysql"
// )

// // ShowAvailableSlots handler to display available vehicles and slots
// func showAvailableSlots(w http.ResponseWriter, r *http.Request) {
// 	// Query to fetch available slots
// 	rows, err := db.Query(`
// 		SELECT vehicleID, AvailableSlotstartDate, AvailableSlotendDate, AvailableSlotstartTime, AvailableSlotendTime
// 		FROM vehicle_schedule
// 		WHERE isAvailable = 1
// 	`)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch data from the database.", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var slots []VehicleSlot
// 	for rows.Next() {
// 		var slot VehicleSlot
// 		if err := rows.Scan(&slot.VehicleID, &slot.AvailableSlotStartDate, &slot.AvailableSlotEndDate, &slot.AvailableSlotStartTime, &slot.AvailableSlotEndTime); err != nil {
// 			http.Error(w, "Unable to process the data.", http.StatusInternalServerError)
// 			return
// 		}
// 		slots = append(slots, slot)
// 	}

// 	// Render HTML template
// 	tmpl := template.Must(template.New("index").Parse(`
// 		<!DOCTYPE html>
// 		<html lang="en">
// 		<head>
// 			<meta charset="UTF-8">
// 			<meta name="viewport" content="width=device-width, initial-scale=1.0">
// 			<title>Car Sharing - Available Slots</title>
// 		</head>
// 		<body>
// 			<h1>Select an Available Slot</h1>
// 			<form action="/reserve" method="POST">
// 				<label for="vehicleID">Vehicle:</label>
// 				<select name="vehicleSlot" id="vehicleID">
// 					{{range .}}
// 						<option value="{{.VehicleID}} - {{.AvailableSlotStartDate}} to {{.AvailableSlotEndDate}} ({{.AvailableSlotStartTime}} - {{.AvailableSlotEndTime}})">{{.VehicleID}} - {{.AvailableSlotStartDate}} to {{.AvailableSlotEndDate}} ({{.AvailableSlotStartTime}} - {{.AvailableSlotEndTime}})</option>
// 					{{end}}
// 				</select><br><br>

// 				<label for="userID">User ID:</label>
// 				<input type="text" name="userID" id="userID" required><br><br>

// 				<input type="submit" value="Reserve Slot">
// 			</form>
// 		</body>
// 		</html>
// 	`))
// 	if err := tmpl.Execute(w, slots); err != nil {
// 		log.Fatal(err)
// 	}
// }

// // ReserveSlot handler to process reservation
// func reserveSlot(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodPost {
// 		// Get form values
// 		vehicleSlot := r.FormValue("vehicleSlot") // E.g., V3 - 2024-12-10 to 2024-12-10 (09:00:00 - 17:30:00)
// 		userID := r.FormValue("userID")

// 		// Parse the vehicleID and time details from the selected slot
// 		var vehicleID, startDate, endDate, startTime, endTime string

// 		// Assuming vehicleSlot format: "V3 - 2024-12-10 to 2024-12-10 (09:00:00 - 17:30:00)"
// 		fmt.Sscanf(vehicleSlot, "%s - %s to %s (%s - %s)", &vehicleID, &startDate, &endDate, &startTime, &endTime)

// 		// Clean up any extraneous characters from startTime and endTime (like closing parentheses)
// 		startTime = strings.TrimRight(startTime, ")")
// 		endTime = strings.TrimRight(endTime, ")")

// 		// Check for missing values
// 		if userID == "" || vehicleID == "" || startDate == "" || endDate == "" || startTime == "" || endTime == "" {
// 			http.Error(w, "Missing userID or vehicleSlot or invalid date/time values", http.StatusBadRequest)
// 			return
// 		}

// 		// Log the input data
// 		log.Printf("Inserting reservation: userID=%s, vehicleID=%s, startDate=%s, endDate=%s, startTime=%s, endTime=%s",
// 			userID, vehicleID, startDate, endDate, startTime, endTime)

// 		// Prepare SQL query to insert into reservation table
// 		stmt, err := db.Prepare(`
//             INSERT INTO reservation (userID, vehicleID, startDate, endDate, startTime, endTime)
//             VALUES (?, ?, ?, ?, ?, ?);
//         `)

// 		if err != nil {
// 			log.Printf("Error preparing statement: %v", err)
// 			http.Error(w, "Error preparing SQL statement", http.StatusInternalServerError)
// 			return
// 		}

// 		// Execute SQL query to insert reservation
// 		_, err = stmt.Exec(userID, vehicleID, startDate, endDate, startTime, endTime)
// 		if err != nil {
// 			log.Printf("Error executing insert: %v", err) // Enhanced error logging
// 			http.Error(w, fmt.Sprintf("Error inserting reservation into database: %v", err), http.StatusInternalServerError)
// 			return
// 		}

//			// Respond back to the user
//			fmt.Fprintf(w, "Reservation for vehicle %s from %s to %s was successfully made.", vehicleID, startDate, endDate)
//		} else {
//			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
//		}
//	}
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type VehicleSlot struct {
	VehicleID              string
	AvailableSlotStartDate string
	AvailableSlotEndDate   string
	AvailableSlotStartTime string
	AvailableSlotEndTime   string
}

func showAvailableSlots(w http.ResponseWriter, r *http.Request) {
	// Query to fetch available slots
	rows, err := db.Query(`
		SELECT vehicleID, AvailableSlotstartDate, AvailableSlotendDate, AvailableSlotstartTime, AvailableSlotendTime
		FROM vehicle_schedule
		WHERE isAvailable = 1
	`)
	if err != nil {
		http.Error(w, "Unable to fetch data from the database.", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var slots []VehicleSlot
	for rows.Next() {
		var slot VehicleSlot
		if err := rows.Scan(&slot.VehicleID, &slot.AvailableSlotStartDate, &slot.AvailableSlotEndDate, &slot.AvailableSlotStartTime, &slot.AvailableSlotEndTime); err != nil {
			http.Error(w, "Unable to process the data.", http.StatusInternalServerError)
			return
		}
		slots = append(slots, slot)
	}

	// Render HTML template
	tmpl := template.Must(template.New("index").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Car Sharing - Available Slots</title>
		</head>
		<body>
			<h1>Select an Available Slot</h1>
			<form action="/reserve" method="POST">
				<label for="vehicleID">Vehicle:</label>
				<select name="vehicleID" id="vehicleID">
					{{range .}}
						<option value="{{.VehicleID}}">{{.VehicleID}} - {{.AvailableSlotStartDate}} to {{.AvailableSlotEndDate}} ({{.AvailableSlotStartTime}} - {{.AvailableSlotEndTime}})</option>
					{{end}}
				</select><br><br>
				<label for="userID">User ID:</label>
				<input type="text" name="userID" id="userID" required><br><br>
				<input type="submit" value="Reserve Slot">
			</form>
		</body>
		</html>
	`))

	if err := tmpl.Execute(w, slots); err != nil {
		log.Fatal(err)
	}
}

func reserveSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		userID := r.FormValue("userID")
		vehicleID := r.FormValue("vehicleID")

		// Check for missing fields
		if userID == "" || vehicleID == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Insert reservation into database
		res, err := db.Exec(`
			INSERT INTO reservation (userID, vehicleID, startDate, endDate, startTime, endTime)
			SELECT ?, vehicleID, AvailableSlotstartDate, AvailableSlotendDate, AvailableSlotstartTime, AvailableSlotendTime
			FROM vehicle_schedule
			WHERE vehicleID = ? AND isAvailable = 1
		`, userID, vehicleID)
		if err != nil {
			http.Error(w, "Error inserting reservation", http.StatusInternalServerError)
			return
		}

		// Get the last inserted reservation ID
		reservationID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error retrieving reservation ID", http.StatusInternalServerError)
			return
		}

		// Redirect to modify booking page
		http.Redirect(w, r, fmt.Sprintf("/modify-booking?reservationID=%d", reservationID), http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
