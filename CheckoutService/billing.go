// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"time"
// )

// type VehicleSlot struct {
// 	VehicleID              string
// 	AvailableSlotStartDate string
// 	AvailableSlotEndDate   string
// 	AvailableSlotStartTime string
// 	AvailableSlotEndTime   string
// 	Duration               string // Store the calculated duration here
// }

// type Member struct {
// 	UserName         string
// 	membershipID     string
// 	typeOfMembership string
// }

// // Function to calculate the duration between the available start and end times
// func calculateSlotDuration(startDate, startTime, endDate, endTime string) (string, error) {
// 	// Combine the date and time to form a full datetime string for start and end
// 	startDateTime := startDate + " " + startTime
// 	endDateTime := endDate + " " + endTime

// 	// Define the datetime format (layout) based on the format of your date and time
// 	const layout = "2006-01-02 15:04:05" // Updated to include seconds

// 	// Parse the start and end date-time strings into time.Time objects
// 	start, err := time.Parse(layout, startDateTime)
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing start date-time: %v", err)
// 	}

// 	end, err := time.Parse(layout, endDateTime)
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing end date-time: %v", err)
// 	}

// 	// Calculate the duration between start and end times
// 	duration := end.Sub(start)

// 	// Convert the duration to days, hours, and minutes
// 	days := int(duration.Hours()) / 24
// 	hours := int(duration.Hours()) % 24
// 	minutes := int(duration.Minutes()) % 60

// 	// Return the formatted duration
// 	return fmt.Sprintf("%d days, %d hours, and %d minutes", days, hours, minutes), nil
// }

// func displayCurrentbooking(w http.ResponseWriter, r *http.Request) {
// 	// Get reservation ID from the query parameters
// 	reservationID := r.URL.Query().Get("reservationID")
// 	if reservationID == "" {
// 		http.Error(w, "Missing reservation ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch current booking details from the reservation table
// 	var userID, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime string
// 	err := db.QueryRow(`
// 		SELECT userID, vehicleID, startDate, endDate, startTime, endTime
// 		FROM reservation
// 		WHERE reservationID = ?`, reservationID).
// 		Scan(&userID, &currentVehicleID, &currentStartDate, &currentEndDate, &currentStartTime, &currentEndTime)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch reservation data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Fetch available slots from the vehicle_schedule table
// 	rows, err := db.Query(`
// 		SELECT vehicleID, AvailableSlotstartDate, AvailableSlotendDate, AvailableSlotstartTime, AvailableSlotendTime
// 		FROM vehicle_schedule
// 		WHERE isAvailable = 1
// 	`)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch available slots", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var slots []VehicleSlot
// 	for rows.Next() {
// 		var slot VehicleSlot
// 		if err := rows.Scan(&slot.VehicleID, &slot.AvailableSlotStartDate, &slot.AvailableSlotEndDate, &slot.AvailableSlotStartTime, &slot.AvailableSlotEndTime); err != nil {
// 			http.Error(w, "Unable to process slot data", http.StatusInternalServerError)
// 			return
// 		}

// 		// Calculate duration for the current slot
// 		duration, err := calculateSlotDuration(slot.AvailableSlotStartDate, slot.AvailableSlotStartTime, slot.AvailableSlotEndDate, slot.AvailableSlotEndTime)
// 		if err != nil {
// 			http.Error(w, "Unable to calculate duration", http.StatusInternalServerError)
// 			return
// 		}

// 		// Add the duration to the slot
// 		slot.Duration = duration

// 		// Append the slot to the slots list
// 		slots = append(slots, slot)
// 	}

// 	// Fetch user details (username and membership) using a JOIN query between users and membership tables
// 	var userName, membershipID, typeOfMembership string
// 	err = db.QueryRow(`
// 		SELECT u.username, u.membershipID, m.typeOfStatus
// 		FROM users u
// 		INNER JOIN membership m ON u.membershipID = m.membershipID
// 		WHERE u.userID = ?`, userID).
// 		Scan(&userName, &membershipID, &typeOfMembership)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch user data", http.StatusInternalServerError)
// 		return
// 	}

// 	// Create a Member struct with the fetched user details
// 	member := Member{
// 		UserName:         userName,
// 		membershipID:     membershipID,
// 		typeOfMembership: typeOfMembership,
// 	}

// 	// Render the HTML template and pass in the data
// 	tmpl := template.Must(template.New("modify").Parse(`
// 		<!DOCTYPE html>
// 		<html lang="en">
// 		<head>
// 			<meta charset="UTF-8">
// 			<meta name="viewport" content="width=device-width, initial-scale=1.0">
// 			<title>Checkout</title>
// 		</head>
// 		<body>
// 			<h1>Reservation Details</h1>
// 			<p><strong>User Name:</strong> {{.UserName}}</p>
// 			<p><strong>Membership ID:</strong> {{.MembershipID}}</p>
// 			<p><strong>Membership Type:</strong> {{.TypeOfMembership}}</p>
// 			<p><strong>Current Booking for Vehicle:</strong> {{.CurrentVehicleID}} from {{.CurrentStartDate}} to {{.CurrentEndDate}} ({{.CurrentStartTime}} - {{.CurrentEndTime}})</p>
// 			<p><strong>Rental Duration:</strong> {{.Duration}}</p>

// 			<h2>Available Vehicle Slots</h2>
// 			<ul>
// 				{{range .Slots}}
// 					<li>{{.VehicleID}}: {{.AvailableSlotStartDate}} - {{.AvailableSlotEndDate}} ({{.AvailableSlotStartTime}} - {{.AvailableSlotEndTime}}) | Duration: {{.Duration}}</li>
// 				{{end}}
// 			</ul>
// 		</body>
// 		</html>
// 	`))

// 	// Pass all the necessary data to the template
// 	err = tmpl.Execute(w, map[string]interface{}{
// 		"UserName":         member.UserName,
// 		"MembershipID":     member.membershipID,
// 		"TypeOfMembership": member.typeOfMembership,
// 		"CurrentVehicleID": currentVehicleID,
// 		"CurrentStartDate": currentStartDate,
// 		"CurrentEndDate":   currentEndDate,
// 		"CurrentStartTime": currentStartTime,
// 		"CurrentEndTime":   currentEndTime,
// 		"Duration":         "",    // You can display the booking duration here if needed
// 		"Slots":            slots, // Passing available slots with calculated durations
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type VehicleSlot struct {
	VehicleID              string
	AvailableSlotStartDate string
	AvailableSlotEndDate   string
	AvailableSlotStartTime string
	AvailableSlotEndTime   string
	Duration               string
}

type Member struct {
	UserName         string
	membershipID     string
	typeOfMembership string
}

// Function to calculate the duration between the available start and end times
func calculateSlotDuration(startDate, startTime, endDate, endTime string) (string, error) {
	// Combine the date and time to form a full datetime string for start and end
	startDateTime := startDate + " " + startTime
	endDateTime := endDate + " " + endTime

	// Define the datetime format (layout) based on the format of your date and time
	const layout = "2006-01-02 15:04:05" // Updated to include seconds

	// Parse the start and end date-time strings into time.Time objects
	start, err := time.Parse(layout, startDateTime)
	if err != nil {
		return "", fmt.Errorf("error parsing start date-time: %v", err)
	}

	end, err := time.Parse(layout, endDateTime)
	if err != nil {
		return "", fmt.Errorf("error parsing end date-time: %v", err)
	}

	// Calculate the duration between start and end times
	duration := end.Sub(start)

	// Convert the duration to days, hours, and minutes
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	// Return the formatted duration
	return fmt.Sprintf("%d days, %d hours, and %d minutes", days, hours, minutes), nil
}

func displayCurrentbooking(w http.ResponseWriter, r *http.Request) {
	// Get reservation ID from the query parameters
	reservationID := r.URL.Query().Get("reservationID")
	if reservationID == "" {
		http.Error(w, "Missing reservation ID", http.StatusBadRequest)
		return
	}

	// Fetch current booking details from the reservation table
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

	// Calculate duration for the current booking
	duration, err := calculateSlotDuration(currentStartDate, currentStartTime, currentEndDate, currentEndTime)
	if err != nil {
		http.Error(w, "Unable to calculate duration", http.StatusInternalServerError)
		return
	}

	// Fetch user details (username and membership) using a JOIN query between users and membership tables
	var userName, membershipID, typeOfMembership string
	err = db.QueryRow(`
		SELECT u.username, u.membershipID, m.typeOfStatus
		FROM users u
		INNER JOIN membership m ON u.membershipID = m.membershipID 
		WHERE u.userID = ?`, userID).
		Scan(&userName, &membershipID, &typeOfMembership)
	if err != nil {
		http.Error(w, "Unable to fetch user data", http.StatusInternalServerError)
		return
	}

	// Create a Member struct with the fetched user details
	member := Member{
		UserName:         userName,
		membershipID:     membershipID,
		typeOfMembership: typeOfMembership,
	}

	// Render the HTML template and pass in the data
	tmpl := template.Must(template.New("modify").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Checkout</title>
		</head>
		<body>
			<h1>Reservation Details</h1>
			<p><strong>User Name:</strong> {{.UserName}}</p>
			<p><strong>Membership ID:</strong> {{.MembershipID}}</p>
			<p><strong>Membership Type:</strong> {{.TypeOfMembership}}</p>
			<p><strong>Current Booking for Vehicle:</strong> {{.CurrentVehicleID}} from {{.CurrentStartDate}} to {{.CurrentEndDate}} ({{.CurrentStartTime}} - {{.CurrentEndTime}})</p>
			<p><strong>Rental Duration:</strong> {{.Duration}}</p>
		</body>
		</html>
	`))

	// Pass all the necessary data to the template
	err = tmpl.Execute(w, map[string]interface{}{
		"UserName":         member.UserName,
		"MembershipID":     member.membershipID,
		"TypeOfMembership": member.typeOfMembership,
		"CurrentVehicleID": currentVehicleID,
		"CurrentStartDate": currentStartDate,
		"CurrentEndDate":   currentEndDate,
		"CurrentStartTime": currentStartTime,
		"CurrentEndTime":   currentEndTime,
		"Duration":         duration, // Showing the duration for the current booking
	})
	if err != nil {
		log.Fatal(err)
	}
}
