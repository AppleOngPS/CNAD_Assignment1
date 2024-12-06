-- Create the CarSharing database and use it
CREATE DATABASE CarSharing;
-- DROP DATABASE carsharing;
USE CarSharing;


-- Create membership table
CREATE TABLE membership (
    membershipID VARCHAR(10) PRIMARY KEY NOT NULL,
    typeOfStatus VARCHAR(10) NOT NULL,
    memberDescriptions VARCHAR(255) NOT NULL,
    discount DECIMAL(5, 2) NOT NULL,
    discountDescription VARCHAR(255) NOT NULL
);


-- Create vehicle table
CREATE TABLE vehicle (
    vehicleID VARCHAR(5) PRIMARY KEY NOT NULL,
    vehicleBrand VARCHAR(255),
   amount DECIMAL(4, 2) NOT NULL
);

-- Create vehicle scheldule
Create table vehicle_schedule (
vehicleID VARCHAR(5) NOT NULL, 
AvailableSlotstartDate DATE NOT NULL,
AvailableSlotendDate DATE NOT NULL,
AvailableSlotstartTime TIME NOT NULL,
AvailableSlotendTime TIME NOT NULL,
isAvailable BOOLEAN NOT NULL DEFAULT true,
FOREIGN KEY (vehicleID) REFERENCES vehicle(vehicleID)
);
-- drop table vehicle_schedule
-- Create promotion table
CREATE TABLE promotion (
    promotionID VARCHAR(5) PRIMARY KEY NOT NULL,
    promotionCode VARCHAR(20) NOT NULL,
    discount DECIMAL(5, 2) NOT NULL,
    description VARCHAR(255) NOT NULL
);

-- Create users table (with auto-increment for userID)
CREATE TABLE users (
    userID INT AUTO_INCREMENT PRIMARY KEY,  -- userID now auto-increments
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(8) NOT NULL,
    membershipID VARCHAR(10) NOT NULL,  
    FOREIGN KEY (membershipID) REFERENCES membership(membershipID)
);

-- Create trackRentalHistory table
CREATE TABLE trackRentalHistory (
    trackRentalHistory VARCHAR(5) PRIMARY KEY NOT NULL,
    userID INT NOT NULL,  -- Changed to INT to match userID type
    vehicleID VARCHAR(5) NOT NULL, 
    vehicleBrand VARCHAR(255) NOT NULL,
    startDate DATE NOT NULL,
    endDate DATE NOT NULL,
    startTime TIME NOT NULL,
    endTime TIME NOT NULL,
    amount DECIMAL(4, 2),
    FOREIGN KEY (userID) REFERENCES users(userID),
    FOREIGN KEY (vehicleID) REFERENCES vehicle(vehicleID)
);



-- Create vehicleStatus table
CREATE TABLE vehicleStatus (
    vehicleID VARCHAR(5) PRIMARY KEY NOT NULL, 
    location VARCHAR(255),
    chargeLevel VARCHAR(255),
    cleanliness VARCHAR(255),
    FOREIGN KEY (vehicleID) REFERENCES vehicle(vehicleID)
);


-- Create reservation table
CREATE TABLE reservation (
    reservationID INT AUTO_INCREMENT PRIMARY KEY,
    userID INT NOT NULL,  -- Changed to INT to match userID type
    vehicleID VARCHAR(5) NOT NULL, 
    startDate DATE NOT NULL,
    endDate DATE NOT NULL,
    startTime TIME NOT NULL,
    endTime TIME NOT NULL,
    FOREIGN KEY (userID) REFERENCES users(userID),
    FOREIGN KEY (vehicleID) REFERENCES vehicle(vehicleID)
);

-- Create billing table
CREATE TABLE billing (
    billingId VARCHAR(5) PRIMARY KEY NOT NULL,
    reservationID INT NOT NULL, 
    amount DECIMAL(4, 2),
    status VARCHAR(10) NOT NULL,
    promotionID VARCHAR(5),  
    FOREIGN KEY (reservationID) REFERENCES reservation(reservationID),
    FOREIGN KEY (promotionID) REFERENCES promotion(promotionID) 
);



-- Insert data into membership table
INSERT INTO membership (membershipID, typeOfStatus, memberdescriptions,discount,discountDescription) VALUES
('M1', 'Basic', 'Basic Membership with limited access',10,'$10 Disoount for all ride'),
('M2', 'Premium', 'Premium Membership with added benefits',20,'$20 Disoount for all ride'),
('M3', 'VIP', 'VIP Membership with all benefits including priority access',30,'$30 Disoount for all ride');

-- Insert data into users table (userID will auto-increment)
INSERT INTO users (username, email, password, membershipID) VALUES
('John', 'john@gmail.com', '12345v6', 'M1'),
('Mary', 'mary@gmail.com', '123g56', 'M2'),
('Mary', 'mary@gmail.com', 'e3456', 'M3'),
('Wong', 'wongSY@gmail.com', '1234v26', 'M2');

-- Insert data into vehicle table
INSERT INTO vehicle (vehicleID, vehicleBrand, amount) VALUES
('V1', 'Honda', 20.00),
('V2', 'Tesla', 15.00),
('V3', 'Nissan ', 18.00);
select*from vehicle;

-- Insert data into promotion table
INSERT INTO promotion (promotionID, promotionCode, discount, description) VALUES
('P001', 'NEWUSER', 20.00, '20% discount for new users'),
('P002', 'HOILDAY', 10.00, '10% discount for holidays');

INSERT INTO vehicle_schedule (vehicleID, availableSlotStartDate, availableSlotEndDate, availableSlotStartTime, availableSlotEndTime, isAvailable)
VALUES
('V1', '2024-12-10', '2024-12-10', '09:10:00', '17:30:00', 1),
('V2', '2024-12-10', '2024-12-10', '09:00:00', '17:00:00', 1),
('V3', '2024-12-10', '2024-12-10', '09:00:00', '17:00:00', 1),
('V3', '2024-12-11', '2024-12-11', '08:00:00', '14:00:00', 1),
('V2', '2024-12-11', '2024-12-11', '10:10:00', '18:45:00', 1),
('V1', '2024-12-12', '2024-12-15', '09:00:00', '16:00:00', 1),
('V2', '2024-12-12', '2024-12-15', '09:00:00', '16:00:00', 1),
('V3', '2024-12-16', '2024-12-16', '07:00:00', '11:00:00', 1),
('V3', '2024-12-16', '2024-12-16', '13:00:00', '17:00:00', 1),
('V1', '2024-12-17', '2024-12-17', '09:00:00', '16:00:00', 1),
('V2', '2024-12-17', '2024-12-17', '09:00:00', '16:00:00', 1),
('V3', '2024-12-17', '2024-12-17', '07:00:00', '11:00:00', 1),
('V3', '2024-12-17', '2024-12-18', '09:50:00', '15:50:00', 1);


-- Insert data into vehicleStatus table
INSERT INTO vehicleStatus (vehicleID, location, chargeLevel, cleanliness) VALUES
('V1', 'Downtown', '80%', 'Clean'),
('V2', 'Uptown', '100%', 'Dirty'),
('V3', 'Suburbs', '50%', 'Clean');

-- Insert data into trackRentalHistory table
INSERT INTO trackRentalHistory (trackRentalHistory, userID, vehicleID, vehicleBrand, startDate, endDate, startTime, endTime, amount) VALUES
('T1', 1, 'V1', 'Honda', '2024-01-10', '2024-01-12', '10:00:00', '14:00:00', 80.00),
('T2', 2, 'V2', 'Tesla', '2024-01-11', '2024-01-13', '11:00:00', '15:00:00', 60.00),
('T3', 3, 'V3', 'Nissan', '2024-01-12', '2024-01-14', '09:00:00', '13:00:00', 72.00),
('T4', 3, 'V3', 'Nissan', '2024-02-03', '2024-02-04', '09:00:00', '13:00:00', 72.00),
('T5', 4, 'V3', 'Nissan', '2024-02-12', '2024-02-14', '09:00:00', '13:00:00', 72.00);

-- Insert data into reservation table
INSERT INTO reservation (reservationID, userID, vehicleID, startDate, endDate, startTime, endTime) VALUES
(1, 1, 'V1', '2024-01-15', '2024-01-16', '10:00:00', '12:00:00'),
(2, 2, 'V2', '2024-01-17', '2024-01-18', '14:00:00', '16:00:00'),
(3, 3, 'V3', '2024-01-19', '2024-01-20', '09:00:00', '11:00:00'),
(4, 4, 'V3', '2024-02-19', '2024-02-20', '09:00:00', '11:00:00');

-- Insert data into billing table
INSERT INTO billing (billingId, reservationID, amount, status, promotionID) VALUES
(1, 1, 40.00, 'Paid','P001'),
(2, 1, 30.00, 'Paid','P002'),
(3, 3, 36.00, 'Pending',NULL),
(4, 3, 36.00, 'Processing','P002');

-- Query data to check insertions
SELECT * FROM billing;
SELECT * FROM reservation;
SELECT * FROM users;
SELECT * FROM membership;
SELECT * FROM trackRentalHistory;
select * from vehicle_schedule;

SELECT * FROM vehicle;
SELECT * FROM vehicleStatus;
SELECT * FROM promotion;

