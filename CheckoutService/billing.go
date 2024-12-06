// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"time"
// )

// // VehicleSlot struct to hold information about available slots
// type VehicleSlot struct {
// 	VehicleID              string
// 	AvailableSlotStartDate string
// 	AvailableSlotEndDate   string
// 	AvailableSlotStartTime string
// 	AvailableSlotEndTime   string
// 	Duration               string
// }

// // Member struct to hold user details
// type Member struct {
// 	UserName         string
// 	membershipID     string
// 	typeOfMembership string
// }

// // Function to calculate the duration between the available start and end times
// func calculateSlotDuration(startDate, startTime, endDate, endTime string) (time.Duration, error) {
// 	startDateTime := startDate + " " + startTime
// 	endDateTime := endDate + " " + endTime

// 	const layout = "2006-01-02 15:04:05" // Date-time format

// 	start, err := time.Parse(layout, startDateTime)
// 	if err != nil {
// 		return 0, fmt.Errorf("error parsing start date-time: %v", err)
// 	}

// 	end, err := time.Parse(layout, endDateTime)
// 	if err != nil {
// 		return 0, fmt.Errorf("error parsing end date-time: %v", err)
// 	}

// 	duration := end.Sub(start) // Calculate the duration

// 	return duration, nil
// }

// // Function to get the membership discount from the database
// func getMembershipDiscount(membershipID string) (float64, error) {
// 	var discount float64
// 	err := db.QueryRow(`
// 		SELECT discount
// 		FROM membership
// 		WHERE membershipID = ?`, membershipID).
// 		Scan(&discount)
// 	if err != nil {
// 		return 0, fmt.Errorf("error fetching discount for membershipID %s: %v", membershipID, err)
// 	}
// 	return discount, nil
// }

// // Function to get the price per hour of the vehicle from the database
// func getVehiclePricePerHour(vehicleID string) (float64, error) {
// 	var pricePerHour float64
// 	err := db.QueryRow(`
// 		SELECT amount
// 		FROM vehicle
// 		WHERE vehicleID = ?`, vehicleID).
// 		Scan(&pricePerHour)
// 	if err != nil {
// 		return 0, fmt.Errorf("error fetching price for vehicleID %s: %v", vehicleID, err)
// 	}
// 	return pricePerHour, nil
// }

// // Function to get available promo codes and their discount
// func getPromoCodes() ([]struct {
// 	Code     string
// 	Discount float64
// }, error) {
// 	rows, err := db.Query(`
// 		SELECT promotionCode, discount
// 		FROM promotion`)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching promo codes: %v", err)
// 	}
// 	defer rows.Close()

// 	var promoCodes []struct {
// 		Code     string
// 		Discount float64
// 	}
// 	for rows.Next() {
// 		var code string
// 		var discount float64
// 		if err := rows.Scan(&code, &discount); err != nil {
// 			return nil, fmt.Errorf("error scanning promo code: %v", err)
// 		}
// 		promoCodes = append(promoCodes, struct {
// 			Code     string
// 			Discount float64
// 		}{Code: code, Discount: discount})
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating over promo codes: %v", err)
// 	}
// 	return promoCodes, nil
// }

// // Function to display current booking details
// func displayCurrentBooking(w http.ResponseWriter, r *http.Request) {
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

// 	// Calculate duration for the current booking
// 	duration, err := calculateSlotDuration(currentStartDate, currentStartTime, currentEndDate, currentEndTime)
// 	if err != nil {
// 		http.Error(w, "Unable to calculate duration", http.StatusInternalServerError)
// 		return
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

// 	// Determine the discount based on membership type (fetch from the database)
// 	discount, err := getMembershipDiscount(membershipID)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch discount", http.StatusInternalServerError)
// 		return
// 	}

// 	// Fetch the price per hour for the vehicle from the database
// 	pricePerHour, err := getVehiclePricePerHour(currentVehicleID)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch price per hour for the vehicle", http.StatusInternalServerError)
// 		return
// 	}

// 	// Calculate the total price based on the rental duration in hours
// 	totalHours := int(duration.Hours()) // Duration in hours
// 	totalPrice := pricePerHour * float64(totalHours)

// 	// Apply the discount to the total price
// 	discountedPrice := totalPrice - discount

// 	// Fetch available promo codes
// 	promoCodes, err := getPromoCodes()
// 	if err != nil {
// 		http.Error(w, "Unable to fetch promo codes", http.StatusInternalServerError)
// 		return
// 	}

// // 	// Render the HTML template and pass in the data
// // 	tmpl := template.Must(template.New("modify").Parse(`
// // 		<!DOCTYPE html>
// // 		<html lang="en">
// // 		<head>
// // 			<meta charset="UTF-8">
// // 			<meta name="viewport" content="width=device-width, initial-scale=1.0">
// // 			<title>Checkout</title>
// // 		</head>
// // 		<body>
// // 			<h1>Reservation Details</h1>
// // 			<p><strong>User Name:</strong> {{.UserName}}</p>
// // 			<p><strong>Membership ID:</strong> {{.MembershipID}}</p>
// // 			<p><strong>Membership Type:</strong> {{.TypeOfMembership}}</p>
// // 			<p><strong>Current Booking for Vehicle:</strong> {{.CurrentVehicleID}} from {{.CurrentStartDate}} to {{.CurrentEndDate}} ({{.CurrentStartTime}} - {{.CurrentEndTime}})</p>
// // 			<p><strong>Rental Duration:</strong> {{.Duration}}</p>
// // 			<p><strong>Discount:</strong> ${{.Discount}}</p>
// // 			<p><strong>Total Price (before discount):</strong> ${{.TotalPrice}}</p>
// // 			<p><strong>Discounted Price:</strong> ${{.DiscountedPrice}}</p>

// // 			<!-- Promo Code Selection -->
// // 			<p><strong>Select Promo Code:</strong></p>
// // 			<select id="promoCodeSelect" onchange="applyPromoCode()">
// // 				<option value="0">Select Promo Code</option>
// // 				{{range .PromoCodes}}
// // 				<option value="{{.Code}}" data-discount="{{.Discount}}">{{.Code}} - ${{.Discount}}</option>
// // 				{{end}}
// // 			</select>
// // 			<p><strong>Final Price after Promo Code:</strong> $<span id="finalPrice">{{.DiscountedPrice}}</span></p>
// // 			<form action="/invoice" method="get">
// //     <input type="hidden" name="reservationID" value="{{.ReservationID}}">
// //     <button type="submit">Go to Invoice</button>
// // </form>

// // </button>
// // 			<script>
// // 				function applyPromoCode() {
// // 					const promoSelect = document.getElementById("promoCodeSelect");
// // 					const discount = promoSelect.selectedOptions[0].getAttribute("data-discount");
// // 					const finalPrice = document.getElementById("finalPrice");
// // 					const discountedPrice = {{.DiscountedPrice}};
// // 					const newPrice = discountedPrice - parseFloat(discount);
// // 					finalPrice.innerText = newPrice.toFixed(2);
// // 				}
// // 			</script>
// // 		</body>
// // 		</html>
// // 	`))

// 	// Pass all the necessary data to the template
// 	err = tmpl.Execute(w, map[string]interface{}{
// 		"UserName":         userName,
// 		"MembershipID":     membershipID,
// 		"TypeOfMembership": typeOfMembership,
// 		"CurrentVehicleID": currentVehicleID,
// 		"CurrentStartDate": currentStartDate,
// 		"CurrentEndDate":   currentEndDate,
// 		"CurrentStartTime": currentStartTime,
// 		"CurrentEndTime":   currentEndTime,
// 		"Duration":         duration,
// 		"Discount":         discount,
// 		"TotalPrice":       totalPrice,
// 		"DiscountedPrice":  discountedPrice,
// 		"PromoCodes":       promoCodes,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// type PromoCode struct {
// 	Code     string
// 	Discount float64
// }

// // Function to display the invoice dynamically based on reservationID
// func displaydetails(w http.ResponseWriter, r *http.Request) {
// 	// Get reservationID from query params
// 	reservationID := r.URL.Query().Get("reservationID")
// 	if reservationID == "" {
// 		http.Error(w, "Missing reservation ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch current booking details from the reservation table based on reservationID
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

// 	// Fetch user details (userName, membershipID, and membershipType) using a JOIN query
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

// 	// Calculate the rental duration
// 	duration, err := calculateSlotDuration(currentStartDate, currentStartTime, currentEndDate, currentEndTime)
// 	if err != nil {
// 		http.Error(w, "Unable to calculate duration", http.StatusInternalServerError)
// 		return
// 	}

// 	// Fetch the price per hour of the vehicle
// 	pricePerHour, err := getVehiclePricePerHour(currentVehicleID)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch price per hour for the vehicle", http.StatusInternalServerError)
// 		return
// 	}

// 	// Calculate the total price based on the rental duration in hours
// 	totalHours := int(duration.Hours()) // Duration in hours
// 	totalPrice := pricePerHour * float64(totalHours)

// 	// Determine the discount based on membership type
// 	discount, err := getMembershipDiscount(membershipID)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch membership discount", http.StatusInternalServerError)
// 		return
// 	}

// 	// Apply membership discount
// 	priceAfterMembershipDiscount := totalPrice - discount

// 	// Fetch promo codes (if any)
// 	promoCodes, err := getPromoCodes()
// 	if err != nil {
// 		http.Error(w, "Unable to fetch promo codes", http.StatusInternalServerError)
// 		return
// 	}

// 	// Apply the first promo code if available
// 	var promoCodeDiscount float64
// 	if len(promoCodes) > 0 {
// 		// Assuming we apply the first promo code
// 		promoCodeDiscount = promoCodes[0].Discount
// 	}

// 	// Final amount to pay after both membership and promo code discounts
// 	finalPrice := priceAfterMembershipDiscount - promoCodeDiscount

// 	// Pass all the fetched data to the template
// 	tmpl := `
// 		<!DOCTYPE html>
// 		<html lang="en">
// 		<head>
// 			<meta charset="UTF-8">
// 			<meta name="viewport" content="width=device-width, initial-scale=1.0">
// 			<title>Invoice</title>
// 		</head>
// 		<body>
// 			<h1>Checkout Page</h1>
// 			<p><strong>User Name:</strong> {{.UserName}}</p>
// 			<p><strong>Membership ID:</strong> {{.MembershipID}}</p>
// 			<p><strong>Membership Type:</strong> {{.TypeOfMembership}}</p>
// 			<p><strong>Current Booking for Vehicle:</strong> {{.CurrentVehicleID}} from {{.CurrentStartDate}} to {{.CurrentEndDate}} ({{.CurrentStartTime}} - {{.CurrentEndTime}})</p>
// 			<p><strong>Rental Duration:</strong> {{.Duration}}</p>
// 			<p><strong>Total Price (before any discount):</strong> ${{.TotalPrice}}</p>
// 			<p><strong>Discount from Membership:</strong> ${{.Discount}}</p>
// 			<p><strong>Promotion Discount:</strong> ${{.PromoCodeDiscount}}</p>
// 			<p><strong>Total Amount (after all discounts):</strong> ${{.FinalPrice}}</p>

// 			<!-- No dropdown for promo code anymore, just final price -->
// 			<p><strong>Final Price:</strong> ${{.FinalPrice}}</p>

// 			<form action="/invoice" method="get">
// 				<input type="hidden" name="reservationID" value="{{.ReservationID}}">
// 				<button type="submit">Go to Invoice</button>
// 			</form>
// 		</body>
// 		</html>
// 	`

// 	// Parse and execute the template
// 	tmplParsed, err := template.New("invoice").Parse(tmpl)
// 	if err != nil {
// 		http.Error(w, "Error parsing template", http.StatusInternalServerError)
// 		return
// 	}

// 	// Pass data to the template
// 	err = tmplParsed.Execute(w, map[string]interface{}{
// 		"UserName":                     userName,
// 		"MembershipID":                 membershipID,
// 		"TypeOfMembership":             typeOfMembership,
// 		"CurrentVehicleID":             currentVehicleID,
// 		"CurrentStartDate":             currentStartDate,
// 		"CurrentEndDate":               currentEndDate,
// 		"CurrentStartTime":             currentStartTime,
// 		"CurrentEndTime":               currentEndTime,
// 		"Duration":                     duration,
// 		"Discount":                     discount,
// 		"TotalPrice":                   totalPrice,
// 		"PriceAfterMembershipDiscount": priceAfterMembershipDiscount,
// 		"PromoCodeDiscount":            promoCodeDiscount,
// 		"FinalPrice":                   finalPrice,
// 		"ReservationID":                reservationID,
// 	})

//		if err != nil {
//			http.Error(w, "Error executing template", http.StatusInternalServerError)
//		}
//	}

// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"net/http"
// 	"net/smtp"
// 	"time"
// )

// // VehicleSlot struct to hold information about available slots
// type VehicleSlot struct {
// 	VehicleID              string
// 	AvailableSlotStartDate string
// 	AvailableSlotEndDate   string
// 	AvailableSlotStartTime string
// 	AvailableSlotEndTime   string
// 	Duration               string
// }

// // Member struct to hold user details
// type Member struct {
// 	UserName         string
// 	membershipID     string
// 	typeOfMembership string
// }

// // Function to calculate the duration between the available start and end times
// func calculateSlotDuration(startDate, startTime, endDate, endTime string) (time.Duration, error) {
// 	startDateTime := startDate + " " + startTime
// 	endDateTime := endDate + " " + endTime

// 	const layout = "2006-01-02 15:04:05" // Date-time format

// 	start, err := time.Parse(layout, startDateTime)
// 	if err != nil {
// 		return 0, fmt.Errorf("error parsing start date-time: %v", err)
// 	}

// 	end, err := time.Parse(layout, endDateTime)
// 	if err != nil {
// 		return 0, fmt.Errorf("error parsing end date-time: %v", err)
// 	}

// 	duration := end.Sub(start) // Calculate the duration

// 	return duration, nil
// }

// // Function to get the membership discount from the database
// func getMembershipDiscount(membershipID string) (float64, error) {
// 	var discount float64
// 	// Sample query, replace with your actual database logic
// 	err := db.QueryRow(`
// 		SELECT discount
// 		FROM membership
// 		WHERE membershipID = ?`, membershipID).
// 		Scan(&discount)
// 	if err != nil {
// 		return 0, fmt.Errorf("error fetching discount for membershipID %s: %v", membershipID, err)
// 	}
// 	return discount, nil
// }

// // Function to get the price per hour of the vehicle from the database
// func getVehiclePricePerHour(vehicleID string) (float64, error) {
// 	var pricePerHour float64
// 	// Sample query, replace with your actual database logic
// 	err := db.QueryRow(`
// 		SELECT amount
// 		FROM vehicle
// 		WHERE vehicleID = ?`, vehicleID).
// 		Scan(&pricePerHour)
// 	if err != nil {
// 		return 0, fmt.Errorf("error fetching price for vehicleID %s: %v", vehicleID, err)
// 	}
// 	return pricePerHour, nil
// }

// // Function to get available promo codes and their discount
// func getPromoCodes() ([]struct {
// 	Code     string
// 	Discount float64
// }, error) {
// 	rows, err := db.Query(`
// 		SELECT promotionCode, discount
// 		FROM promotion`)
// 	if err != nil {
// 		return nil, fmt.Errorf("error fetching promo codes: %v", err)
// 	}
// 	defer rows.Close()

// 	var promoCodes []struct {
// 		Code     string
// 		Discount float64
// 	}
// 	for rows.Next() {
// 		var code string
// 		var discount float64
// 		if err := rows.Scan(&code, &discount); err != nil {
// 			return nil, fmt.Errorf("error scanning promo code: %v", err)
// 		}
// 		promoCodes = append(promoCodes, struct {
// 			Code     string
// 			Discount float64
// 		}{Code: code, Discount: discount})
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating over promo codes: %v", err)
// 	}
// 	return promoCodes, nil
// }

// // Function to send the reservation details via email
// func sendEmail(userEmail, userName, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime string, totalPrice, discountedPrice float64, reservationID string) error {
// 	// Sender email credentials
// 	from := "ongapple1@gmail.com"
// 	password := "dsvw cwzk oieg nglb"

// 	// SMTP server configuration
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"

// 	// Construct the URL for the invoice page
// 	invoiceURL := fmt.Sprintf("http://localhost:8083/invoice?reservationID=%s", reservationID)

// 	// Compose the email content
// 	subject := "Your Vehicle Reservation Details"
// 	body := fmt.Sprintf(`
// 		Hello %s,

// 		Thank you for your reservation! Here are the details:

// 		Vehicle ID: %s
// 		Reservation Start Date: %s
// 		Reservation End Date: %s
// 		Reservation Start Time: %s
// 		Reservation End Time: %s

// 		Total Price: $%.2f
// 		Discounted Price: $%.2f

// 		You can view and download your invoice here: %s

// 		Thank you for choosing us!
// 		`, userName, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime, totalPrice, discountedPrice, invoiceURL)

// 	// Set up the message
// 	msg := []byte("To: " + userEmail + "\r\n" +
// 		"Subject: " + subject + "\r\n" +
// 		"\r\n" +
// 		body)

// 	// Authenticate with the SMTP server
// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	// Send the email
// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{userEmail}, msg)
// 	if err != nil {
// 		return fmt.Errorf("failed to send email: %v", err)
// 	}

// 	return nil
// }

// // Function to display current booking details
// func displayCurrentBooking(w http.ResponseWriter, r *http.Request) {
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

// 	// Calculate duration for the current booking
// 	duration, err := calculateSlotDuration(currentStartDate, currentStartTime, currentEndDate, currentEndTime)
// 	if err != nil {
// 		http.Error(w, "Unable to calculate duration", http.StatusInternalServerError)
// 		return
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

// 	// Determine the discount based on membership type (fetch from the database)
// 	discount, err := getMembershipDiscount(membershipID)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch discount", http.StatusInternalServerError)
// 		return
// 	}

// 	// Fetch the price per hour for the vehicle from the database
// 	pricePerHour, err := getVehiclePricePerHour(currentVehicleID)
// 	if err != nil {
// 		http.Error(w, "Unable to fetch price per hour for the vehicle", http.StatusInternalServerError)
// 		return
// 	}

// 	// Calculate the total price based on the rental duration in hours
// 	totalHours := int(duration.Hours()) // Duration in hours
// 	totalPrice := pricePerHour * float64(totalHours)

// 	// Apply the discount to the total price
// 	discountedPrice := totalPrice - discount

// 	// Fetch available promo codes
// 	promoCodes, err := getPromoCodes()
// 	if err != nil {
// 		http.Error(w, "Unable to fetch promo codes", http.StatusInternalServerError)
// 		return
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
// 			<p><strong>Discount:</strong> ${{.Discount}}</p>
// 			<p><strong>Total Price (before discount):</strong> ${{.TotalPrice}}</p>
// 			<p><strong>Discounted Price:</strong> ${{.DiscountedPrice}}</p>

// 			<!-- Promo Code Selection -->
// 			<p><strong>Select Promo Code:</strong></p>
// 			<select id="promoCodeSelect" onchange="applyPromoCode()">
// 				<option value="0">Select Promo Code</option>
// 				{{range .PromoCodes}}
// 				<option value="{{.Code}}" data-discount="{{.Discount}}">{{.Code}} - ${{.Discount}}</option>
// 				{{end}}
// 			</select>
// 			<p><strong>Final Price after Promo Code:</strong> $<span id="finalPrice">{{.DiscountedPrice}}</span></p>
// 			<form action="/invoice" method="get">
//     <input type="hidden" name="reservationID" value="{{.ReservationID}}">
//     <button type="submit">Go to Invoice</button>
// </form>

// </button>
// 			<script>
// 				function applyPromoCode() {
// 					const promoSelect = document.getElementById("promoCodeSelect");
// 					const discount = promoSelect.selectedOptions[0].getAttribute("data-discount");
// 					const finalPrice = document.getElementById("finalPrice");
// 					const discountedPrice = {{.DiscountedPrice}};
// 					const newPrice = discountedPrice - parseFloat(discount);
// 					finalPrice.innerText = newPrice.toFixed(2);
// 				}
// 			</script>
// 		</body>
// 		</html>
// 	`))

// 	err = tmpl.Execute(w, struct {
// 		UserName         string
// 		MembershipID     string
// 		TypeOfMembership string
// 		CurrentVehicleID string
// 		CurrentStartDate string
// 		CurrentEndDate   string
// 		CurrentStartTime string
// 		CurrentEndTime   string
// 		Duration         string
// 		Discount         float64
// 		TotalPrice       float64
// 		DiscountedPrice  float64
// 		PromoCodes       []struct {
// 			Code     string
// 			Discount float64
// 		}
// 	}{
// 		UserName:         userName,
// 		MembershipID:     membershipID,
// 		TypeOfMembership: typeOfMembership,
// 		CurrentVehicleID: currentVehicleID,
// 		CurrentStartDate: currentStartDate,
// 		CurrentEndDate:   currentEndDate,
// 		CurrentStartTime: currentStartTime,
// 		CurrentEndTime:   currentEndTime,
// 		Duration:         fmt.Sprintf("%.2f hours", duration.Hours()),
// 		Discount:         discount,
// 		TotalPrice:       totalPrice,
// 		DiscountedPrice:  discountedPrice,
// 		PromoCodes:       promoCodes,
// 	})

// 	if err != nil {
// 		http.Error(w, "Unable to render template", http.StatusInternalServerError)
// 		return
// 	}

// 	// Optionally, send email notification with invoice link
// 	err = sendEmail("user@example.com", userName, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime, totalPrice, discountedPrice, reservationID)
// 	if err != nil {
// 		http.Error(w, "Unable to send email", http.StatusInternalServerError)
// 		return
// 	}
// }

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
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
	// Sample query, replace with your actual database logic
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
	// Sample query, replace with your actual database logic
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

// Function to get available promo codes and their discount
func getPromoCodes() ([]struct {
	Code     string
	Discount float64
}, error) {
	rows, err := db.Query(`
		SELECT promotionCode, discount
		FROM promotion`)
	if err != nil {
		return nil, fmt.Errorf("error fetching promo codes: %v", err)
	}
	defer rows.Close()

	var promoCodes []struct {
		Code     string
		Discount float64
	}
	for rows.Next() {
		var code string
		var discount float64
		if err := rows.Scan(&code, &discount); err != nil {
			return nil, fmt.Errorf("error scanning promo code: %v", err)
		}
		promoCodes = append(promoCodes, struct {
			Code     string
			Discount float64
		}{Code: code, Discount: discount})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over promo codes: %v", err)
	}
	return promoCodes, nil
}

// Function to send the reservation details via email
func sendEmail(userEmail, userName, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime string, totalPrice, discountedPrice float64, reservationID string) error {
	// Sender email credentials
	from := "ongapple1@gmail.com"
	password := "dsvw cwzk oieg nglb"

	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Construct the URL for the invoice page
	invoiceURL := fmt.Sprintf("http://localhost:8083/invoice?reservationID=%s", reservationID)

	// Compose the email content
	subject := "Your Vehicle Reservation Details"
	body := fmt.Sprintf(`Hello %s,

Thank you for your reservation! Here are the details:

Vehicle ID: %s
Reservation Start Date: %s
Reservation End Date: %s
Reservation Start Time: %s
Reservation End Time: %s

Total Price: $%.2f
Discounted Price: $%.2f

You can view and download your invoice here: %s

Thank you for choosing us!
`, userName, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime, totalPrice, discountedPrice, invoiceURL)

	// Set up the message
	msg := []byte("To: " + userEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{userEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// Function to handle reservation confirmation
func confirmReservation(w http.ResponseWriter, r *http.Request) {
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

	// Fetch user details (username and membership) using a JOIN query between users and membership tables
	var userName, membershipID, typeOfMembership, userEmail string
	err = db.QueryRow(`
		SELECT u.username, u.membershipID, m.typeOfStatus, u.email
		FROM users u
		INNER JOIN membership m ON u.membershipID = m.membershipID
		WHERE u.userID = ?`, userID).
		Scan(&userName, &membershipID, &typeOfMembership, &userEmail)
	if err != nil {
		http.Error(w, "Unable to fetch user data", http.StatusInternalServerError)
		return
	}

	// Calculate duration for the current booking
	duration, err := calculateSlotDuration(currentStartDate, currentStartTime, currentEndDate, currentEndTime)
	if err != nil {
		http.Error(w, "Unable to calculate duration", http.StatusInternalServerError)
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

	// Fetch available promo codes
	promoCodes, err := getPromoCodes()
	if err != nil {
		http.Error(w, "Unable to fetch promo codes", http.StatusInternalServerError)
		return
	}

	// Check if a valid promo code was used
	promoCode := r.URL.Query().Get("promoCode") // Assume promo code is passed as a query parameter
	for _, code := range promoCodes {
		if code.Code == promoCode {
			discountedPrice -= code.Discount // Apply the promo code discount
			break
		}
	}

	// Send the confirmation email to the user with the reservation details
	err = sendEmail(userEmail, userName, currentVehicleID, currentStartDate, currentEndDate, currentStartTime, currentEndTime, totalPrice, discountedPrice, reservationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send email: %v", err), http.StatusInternalServerError)
		return
	}

	// Optionally, render a confirmation page or return a response
	fmt.Fprintf(w, "Reservation confirmed! A confirmation email has been sent to %s", userEmail)
}

// Display Current Booking
func displayCurrentBooking(w http.ResponseWriter, r *http.Request) {
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

	// Use templates to render the page
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

			<!-- Promo Code Selection -->
			<p><strong>Select Promo Code:</strong></p>
			<select id="promoCodeSelect" onchange="applyPromoCode()">
				<option value="0">Select Promo Code</option>
				{{range .PromoCodes}}
				<option value="{{.Code}}" data-discount="{{.Discount}}">{{.Code}} - ${{.Discount}}</option>
				{{end}}
			</select>
			<p><strong>Final Price after Promo Code:</strong> $<span id="finalPrice">{{.DiscountedPrice}}</span></p>
			<form action="/invoice" method="get">
    <input type="hidden" name="reservationID" value="{{.ReservationID}}">
    <button type="submit">Go to Invoice</button>
</form>

</button>
			<script>
				function applyPromoCode() {
					const promoSelect = document.getElementById("promoCodeSelect");
					const discount = promoSelect.selectedOptions[0].getAttribute("data-discount");
					const finalPrice = document.getElementById("finalPrice");
					const discountedPrice = {{.DiscountedPrice}};
					const newPrice = discountedPrice - parseFloat(discount);
					finalPrice.innerText = newPrice.toFixed(2);
				}
			</script>
		</body>
		</html>
	`))

	// Render the template with current booking data
	tmpl.Execute(w, struct {
		CurrentVehicleID string
		CurrentStartDate string
		CurrentEndDate   string
		CurrentStartTime string
		CurrentEndTime   string
	}{
		CurrentVehicleID: currentVehicleID,
		CurrentStartDate: currentStartDate,
		CurrentEndDate:   currentEndDate,
		CurrentStartTime: currentStartTime,
		CurrentEndTime:   currentEndTime,
	})
}
