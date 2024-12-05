package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// VehicleSlot struct to hold information about available slots
type VehicleSlot struct {
	VehicleID              string
	AvailableSlotStartDate string
	AvailableSlotEndDate   string
	AvailableSlotStartTime string
	AvailableSlotEndTime   string
	Duration               string
}

// Member struct to hold user details
type Member struct {
	UserName         string
	membershipID     string
	typeOfMembership string
}

// Function to calculate the duration between the available start and end times
func calculateSlotDuration(startDate, startTime, endDate, endTime string) (time.Duration, error) {
	startDateTime := startDate + " " + startTime
	endDateTime := endDate + " " + endTime

	const layout = "2006-01-02 15:04:05" // Date-time format

	start, err := time.Parse(layout, startDateTime)
	if err != nil {
		return 0, fmt.Errorf("error parsing start date-time: %v", err)
	}

	end, err := time.Parse(layout, endDateTime)
	if err != nil {
		return 0, fmt.Errorf("error parsing end date-time: %v", err)
	}

	duration := end.Sub(start) // Calculate the duration

	return duration, nil
}

// Function to get the membership discount from the database
func getMembershipDiscount(membershipID string) (float64, error) {
	var discount float64
	err := db.QueryRow(`
		SELECT discount
		FROM membership
		WHERE membershipID = ?`, membershipID).
		Scan(&discount)
	if err != nil {
		return 0, fmt.Errorf("error fetching discount for membershipID %s: %v", membershipID, err)
	}
	return discount, nil
}

// Function to get the price per hour of the vehicle from the database
func getVehiclePricePerHour(vehicleID string) (float64, error) {
	var pricePerHour float64
	err := db.QueryRow(`
		SELECT amount
		FROM vehicle
		WHERE vehicleID = ?`, vehicleID).
		Scan(&pricePerHour)
	if err != nil {
		return 0, fmt.Errorf("error fetching price for vehicleID %s: %v", vehicleID, err)
	}
	return pricePerHour, nil
}

// Function to display current booking details
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

	// Determine the discount based on membership type (fetch from the database)
	discount, err := getMembershipDiscount(membershipID)
	if err != nil {
		http.Error(w, "Unable to fetch discount", http.StatusInternalServerError)
		return
	}

	// Fetch the price per hour for the vehicle from the database
	pricePerHour, err := getVehiclePricePerHour(currentVehicleID)
	if err != nil {
		http.Error(w, "Unable to fetch price per hour for the vehicle", http.StatusInternalServerError)
		return
	}

	// Calculate the total price based on the rental duration in hours
	totalHours := int(duration.Hours()) // Duration in hours
	totalPrice := pricePerHour * float64(totalHours)

	// Apply the discount to the total price
	discountedPrice := totalPrice - discount

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
			<p><strong>Discount:</strong> ${{.Discount}}</p>
			<p><strong>Total Price (before discount):</strong> ${{.TotalPrice}}</p>
			<p><strong>Discounted Price:</strong> ${{.DiscountedPrice}}</p>
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
		"Duration":         duration,
		"Discount":         discount,
		"TotalPrice":       totalPrice,
		"DiscountedPrice":  discountedPrice,
	})
	if err != nil {
		log.Fatal(err)
	}
}
