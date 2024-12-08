package main

import (
	"html/template"
	"log"
	"net/http"
)

// VehicleSlot represents available slots in the vehicle schedule

func showModifyBooking(w http.ResponseWriter, r *http.Request) {
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
        <title>Modify Current Booking</title>
    </head>
    <body>
        <h1>Modify Your Booking</h1>
        <p>User ID: {{.UserID}}</p>
        <p>Current Booking for Vehicle: {{.CurrentVehicleID}} from {{.CurrentStartDate}} to {{.CurrentEndDate}} ({{.CurrentStartTime}} - {{.CurrentEndTime}})</p>
        
        <!-- Update Booking drop down list -->
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
        
        <!-- Delete Booking Form -->
        <form action="/delete-booking" method="POST" style="margin-top: 20px;">
            <input type="hidden" name="reservationID" value="{{.ReservationID}}">
            <button type="submit">Delete Booking</button>
        </form>

        <form action="http://localhost:8083/displayCurrentBooking" method="GET">
            <input type="submit" value="Check Out">
        </form>
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

func modifyBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get form values
		reservationID := r.FormValue("reservationID")
		vehicleID := r.FormValue("vehicleID")

		// Check for missing values
		if reservationID == "" || vehicleID == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Update reservation in the database
		_, err := db.Exec(`
			UPDATE reservation 
			SET vehicleID = ? 
			WHERE reservationID = ?`, vehicleID, reservationID)
		if err != nil {
			http.Error(w, "Error updating reservation", http.StatusInternalServerError)
			return
		}

		// Redirect to confirmation page or back to slots
		//http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func deleteBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Get reservation ID from the form
		reservationID := r.FormValue("reservationID")

		// Check if reservation ID is missing
		if reservationID == "" {
			http.Error(w, "Missing reservation ID", http.StatusBadRequest)
			return
		}

		// Delete the reservation from the database
		_, err := db.Exec(`
			DELETE FROM reservation 
			WHERE reservationID = ?`, reservationID)
		if err != nil {
			http.Error(w, "Error deleting reservation", http.StatusInternalServerError)
			return
		}

		// Redirect to the homepage or another page
		//http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
