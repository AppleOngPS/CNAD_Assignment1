package main

import (
	"net/http"
)

// Handle the reservation page
func Reservation(w http.ResponseWriter, r *http.Request) {
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

				<label for="vehicleID">Vehicle ID:</label>
				<input type="text" id="vehicleID" name="vehicleID" required><br><br>

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
	w.Write([]byte(tmpl))
}

// Handle the reservation form submission
func reservationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		reservationID := r.FormValue("reservationID")
		userID := r.FormValue("userID")
		vehicleID := r.FormValue("vehicleID")
		startDate := r.FormValue("startDate")
		endDate := r.FormValue("endDate")
		startTime := r.FormValue("startTime")
		endTime := r.FormValue("endTime")

		// Insert the new reservation into the database
		_, err := db.Exec(
			"INSERT INTO reservation (reservationID, userID, vehicleID, startDate, endDate, startTime, endTime) VALUES (?, ?, ?, ?, ?, ?, ?)",
			reservationID, userID, vehicleID, startDate, endDate, startTime, endTime,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Redirect to a success page or display a success message
		http.Redirect(w, r, "/reservation/success", http.StatusSeeOther)
	}
}

// Handle the reservation success page
func reservationSuccess(w http.ResponseWriter, r *http.Request) {
	tmpl := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Reservation Success</title>
		</head>
		<body>
			<h2>Reservation Successful!</h2>
			<p>Your car reservation has been submitted successfully.</p>
		</body>
		</html>
	`
	w.Write([]byte(tmpl))
}
