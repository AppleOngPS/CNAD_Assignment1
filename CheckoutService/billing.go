package main

import (
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

func displayCurrentbooking(w http.ResponseWriter, r *http.Request) {
	// Get reservation ID
	reservationID := r.URL.Query().Get("reservationID")
	if reservationID == "" {
		http.Error(w, "Missing reservation ID", http.StatusBadRequest)
		return
	}

	// Fetch current booking details
	var userID, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime string
	err := db.QueryRow(`
		SELECT userID, vehicleID, startDate, endDate, startTime, endTime 
		FROM reservation 
		WHERE reservationID = ?`, reservationID).
		Scan(&userID, &currentVehicleID, &currentStartDate, &currentEndDate, &currentStartTime, &currentEndTime)
	if err != nil {
		http.Error(w, "Unable to fetch reservation data", http.StatusInternalServerError)
		return
	}

	// Fetch available slots
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

	// Render modify booking page
	tmpl := template.Must(template.New("modify").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Checkout</title>
		</head>
		<body>
			<h1>Reservation details</h1>
			<p>User ID: {{.UserID}}</p>
			<p>Current Booking for Vehicle: {{.CurrentVehicleID}} from {{.CurrentStartDate}} to {{.CurrentEndDate}} ({{.CurrentStartTime}} - {{.CurrentEndTime}})</p>
			
		
		</body>
		</html>
	`))

	err = tmpl.Execute(w, map[string]interface{}{
		"ReservationID":    reservationID,
		"UserID":           userID,
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
