package main

import (
	"net/http"
)

// Handle the reservation page
func Reservation(w http.ResponseWriter, r *http.Request) {
	// Create the HTML form
	tmpl := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Car Reservation</title>
		</head>
		<body>
			<h2>Car Reservation Form</h2>
			<form action="/reservation/submit" method="POST">
				<label for="reservationID">Reservation ID:</label>
				<input type="text" id="reservationID" name="reservationID" required><br><br>

				<label for="userID">User ID:</label>
				<input type="number" id="userID" name="userID" required><br><br>

				<label for="startDate">Start Date:</label>
				<input type="date" id="startDate" name="startDate" required><br><br>

				<label for="endDate">End Date:</label>
				<input type="date" id="endDate" name="endDate" required><br><br>

				<label for="startTime">Start Time:</label>
				<input type="time" id="startTime" name="startTime" required><br><br>

				<label for="endTime">End Time:</label>
				<input type="time" id="endTime" name="endTime" required><br><br>

				<input type="submit" value="Reserve">
			</form>
		</body>
		</html>
	`

	// Write the form to the response
	w.Write([]byte(tmpl))
}
