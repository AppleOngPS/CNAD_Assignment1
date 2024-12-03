package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// Handle the car listing page
func CarListingHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Query to get all car listings from the vehicle table
	rows, err := db.Query(`
		SELECT 
			vehicleID,
			vehicleBrand,
			startDate,
			endDate,
			startTime,
			endTime,
			amount
		FROM vehicle
	`)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Error retrieving car listings", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Slice to hold all the cars
	var cars []Car

	// Loop through the rows and populate the car slice
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.VehicleID, &car.VehicleBrand, &car.StartDate, &car.EndDate, &car.StartTime, &car.EndTime, &car.Amount); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		cars = append(cars, car)
	}

	// Check for errors after iterating over the rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	// Create the HTML response manually
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTML layout
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Car Listing</title>
		</head>
		<body>
			<h1>Available Cars for Rent</h1>
			<div>
	`

	// Loop through each car and append the HTML for each
	for _, car := range cars {
		htmlContent += fmt.Sprintf(`
			<div>
				<h2>%s</h2>
				<p><strong>Vehicle ID:</strong> %s</p>
				<p><strong>Available From:</strong> %s To %s</p>
				<p><strong>Available Time:</strong> %s - %s</p>
				<p><strong>Price per Hour:</strong> $%.2f</p>
				<!-- General Reserve Button -->
				<form action="/reservation" method="GET">
					<input type="submit" value="Reserve">
				</form>
			</div>`, car.VehicleBrand, car.VehicleID, car.StartDate, car.EndDate, car.StartTime, car.EndTime, car.Amount)
	}

	htmlContent += `
			</div>
		</body>
		</html>`

	// Write the final HTML content to the response
	w.Write([]byte(htmlContent))
}
