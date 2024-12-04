package main

import (
	"html/template"
	"log"
	"net/http"
)

// VehicleSlot represents a single slot in the schedule
type VehicleSlot struct {
	VehicleID              string
	AvailableSlotStartDate string
	AvailableSlotEndDate   string
	AvailableSlotStartTime string
	AvailableSlotEndTime   string
}

// showAvailableSlots: Show available vehicle slots

// showModifyBookingHandler: Show the form to modify an existing booking
func showModifyBookingHandler(w http.ResponseWriter, r *http.Request) {
	// Get the reservation ID from the query parameters
	reservationID := r.URL.Query().Get("reservationID")
	if reservationID == "" {
		http.Error(w, "Missing reservation ID", http.StatusBadRequest)
		return
	}

	// Query the reservation details
	var currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime string
	err := db.QueryRow(`SELECT vehicleID, startDate, endDate, startTime, endTime FROM reservation WHERE reservationID = ?`, reservationID).
		Scan(&currentVehicleID, &currentStartDate, &currentEndDate, &currentStartTime, &currentEndTime)
	if err != nil {
		http.Error(w, "Unable to fetch reservation data", http.StatusInternalServerError)
		return
	}

	// Query to fetch available slots
	rows, err := db.Query(`
		SELECT vehicleID, AvailableSlotstartDate, AvailableSlotendDate, AvailableSlotstartTime, AvailableSlotendTime
		FROM vehicle_schedule
		WHERE isAvailable = 1
	`)
	if err != nil {
		http.Error(w, "Unable to fetch available slots", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var slots []VehicleSlot
	for rows.Next() {
		var slot VehicleSlot
		if err := rows.Scan(&slot.VehicleID, &slot.AvailableSlotStartDate, &slot.AvailableSlotEndDate, &slot.AvailableSlotStartTime, &slot.AvailableSlotEndTime); err != nil {
			http.Error(w, "Unable to process slot data", http.StatusInternalServerError)
			return
		}
		slots = append(slots, slot)
	}

	// Render the form to modify the booking
	tmpl := template.Must(template.New("modify").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Modify Your Booking</title>
		</head>
		<body>
			<h1>Modify Your Booking</h1>
			<p>Current Booking for Vehicle: {{.CurrentVehicleID}} from {{.CurrentStartDate}} to {{.CurrentEndDate}} ({{.CurrentStartTime}} - {{.CurrentEndTime}})</p>
			<form action="/update-booking" method="POST">
				<input type="hidden" name="reservationID" value="{{.ReservationID}}">
				<label for="vehicleID">Select New Vehicle Slot:</label>
				<select name="vehicleID" id="vehicleID">
					{{range .Slots}}
						<option value="{{.VehicleID}}">{{.VehicleID}} - {{.AvailableSlotStartDate}} to {{.AvailableSlotEndDate}} ({{.AvailableSlotStartTime}} - {{.AvailableSlotEndTime}})</option>
					{{end}}
				</select><br><br>
				<input type="submit" value="Update Booking">
			</form>
		</body>
		</html>
	`))

	err = tmpl.Execute(w, map[string]interface{}{
		"ReservationID":    reservationID,
		"CurrentVehicleID": currentVehicleID,
		"CurrentStartDate": currentStartDate,
		"CurrentEndDate":   currentEndDate,
		"CurrentStartTime": currentStartTime,
		"CurrentEndTime":   currentEndTime,
		"Slots":            slots,
	})
	if err != nil {
		log.Fatal(err)
	}
}
