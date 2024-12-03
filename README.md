# CNAD_Assignment1


# Design consideration of your microservices

**Service Decomposition:**<br /> 

**User Management**<br /> 
**Login**<br /> 
Purpose: Allow user to login to view using their own credential<br /> 
Backend: Verify the userID and password against the database, if it exist it will proceed to the car listing page<br /> 
SQL: table check if there is data exist and check if they enter the details correctly <br /> 

**Signup** <br /> 
Purpose: Allow user to login to create their own credential<br /> 
Backend:When user enter new record it will insert new row of data into the table<br /> 
SQL: check if they enter the details correctly <br /> 

**Profile**<br /> 
Purpose: Manage user profile , view past booking and manage existing booking<br /> 
Backend:When user modify the details the database will update or delete the record in the database<br /> 
SQL: check if they enter the details correctly <br /> 

**View Listing Service**
**Car Listing**<br /> 
Purpose: View all listing of car<br /> 
Backend: it will retreve the data from the database and display in the website
SQL: retreive the data from the database

**Reservation service**
Reservation: Create reservation<br /> 
Purpose: Create a new booking<br /> 
Backend: It will insert a new record when user create a new booking
SQL: insert a new record

**Checkout Service**
Checkout: calculate amount and create receipt<br /> 
purpose: To create a invoice base on what the user had reserved
Backend: to generate a invoice to allow user to see
SQL: Insert a new record

**Others**
**Inter-Service communication**<br /> 
User HTTP/REST for synchronous call (eg from car listing to reservation)<br /> 

**Security**<br /> 
Secure all sensetive data such as user password using JWT Token<br /> 

**Database**
All microservices link to a single database to avoid data duplication and ensure data consistency. However, denormalization may be used in some cases to improve query performance, depending on the service's needs.

# Architecture diagram

![architecture diagram](https://github.com/user-attachments/assets/cf4189b9-bc90-41b2-9a2f-860d8ac18e3f)

# Instructions for setting up and running your microservices





# Reference
https://www.designgurus.io/answers/detail/what-are-the-key-considerations-for-designing-a-microservices-architecture
