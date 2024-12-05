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
			<h1>Reserve a slot</h1>
			<h2>Select a slot that you want!</h2>
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
			LIMIT 1;
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
