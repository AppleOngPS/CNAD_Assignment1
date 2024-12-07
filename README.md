# CNAD_Assignment1


# Design consideration of your microservices<br /> <br /> 

**Service Decomposition:**<br /> 

**User Management**<br /> 
1. Login <br /> 
Purpose: Allow user to login to view using their own credential<br /> 
Backend: Verify the userID and password against the database, if it exist it will proceed to the car listing page<br /> 
SQL: table check if there is data exist and check if they enter the details correctly <br /> 

2. Signup <br /> 
Purpose: Allow user to login to create their own credential<br /> 
Backend:When user enter new record it will insert new row of data into the table<br /> 
SQL: check if they enter the details correctly <br /> 

3. Profile <br /> 
Purpose: Manage user profile , view past booking and manage existing booking<br /> 
Backend:When user modify the details the database will update or delete the record in the database<br /> 
SQL: check if they enter the details correctly <br /> 

**View Listing Service**
1. Car Listing <br /> 
Purpose: View all listing of car<br /> 
Backend: it will retreve the data from the database and display in the website
SQL: retreive the data from the database

**Reservation service**
1. Reservation: Create reservation<br /> 
Purpose: Create a new booking<br /> 
Backend: It will insert a new record when user create a new booking
SQL: insert a new record

**Checkout Service**
1. Checkout: calculate amount and create receipt<br /> 
purpose: To create a invoice base on what the user had reserved
Backend: to generate a invoice to allow user to see
SQL: Insert a new record

**Others**
**Inter-Service communication**<br /> 
User HTTP/REST for synchronous call (eg from car listing to reservation)<br /> 

**Security**<br /> 
Secure all sensetive data such as user password using bcrypt<br /> 
Verification email sent when signing up a new account <br /> 
Verification email sent when complete the payment and will sent the invoice link <br /> 

**Database**
All microservices link to a single database to avoid data duplication and ensure data consistency. However, denormalization may be used in some cases to improve query performance, depending on the service's needs.

# Architecture diagram
<br /> 

![architecture diagram](https://github.com/user-attachments/assets/1a8febc5-4675-494d-989e-ceacd7bc7ce8)  <br /> 

# Instructions for setting up and running your microservices<br /> <br /> 
<br /> 
1. Open 4 different terminal and run the service (follow below screenshot)<br /> <br /> 
   
For UserService<br /> 
![UserService_Terminal_ScreenShot](https://github.com/user-attachments/assets/a2012b8b-a2c2-442c-a059-8097cceedd6f)<br /> 

For ViewCarListing<br /> 
![Carlisting_Terminal_Screenshot](https://github.com/user-attachments/assets/e07710c0-8db0-40e3-b2c5-5101ce6c4bff)<br /> 

For ReservationService<br /> 
![ReservationService_Terminal_Screenshot](https://github.com/user-attachments/assets/2b62baec-d696-4b25-a567-b881b17cf716)<br /> 

For CheckoutService
![CheckoutService_Terminal_Screenshot](https://github.com/user-attachments/assets/d52639d7-aa0c-456a-a75c-87fb5881b209)

2. To access the page for each service (refer to the routing in each main file)<br /> <br /> 

To go to signup page: <br /> 
![image](https://github.com/user-attachments/assets/a2e53696-d5bd-4c5e-b24c-49504669ea4c)<br /> 

To go to login page:<br /> 
![image](https://github.com/user-attachments/assets/c68360cc-f072-4159-852f-83e2a776fba1)<br /> 

To go to profile page<br /> 
![image](https://github.com/user-attachments/assets/77110b2c-11ee-4d63-8f9b-c68c372acdf8)<br /> 


To go to car listing page<br /> 
![image](https://github.com/user-attachments/assets/c7cefcff-20e9-4e60-9203-39e2c4d07539)<br /> 

To go to reserve slot <br /> 
![image](https://github.com/user-attachments/assets/0ea33f75-18f7-49ec-a8c2-56b9272083b9)<br /> 

To modify booking<br /> 
![image](https://github.com/user-attachments/assets/f247b0ca-138d-4d18-b8b2-5fc6131700d4)<br /> 

To display current booking to pay<br /> 
![image](https://github.com/user-attachments/assets/6b546679-a28d-446b-ab49-db4eb8098e91)<br /> 

To send email confirmation <br /> 
![image](https://github.com/user-attachments/assets/3b3ea8de-9529-46b0-bfe3-f2c9416c642d)<br /> 

Generate Invoice<br /> 
![image](https://github.com/user-attachments/assets/b2875c96-60dd-42e5-86f7-832d136bd1c4)<br /> 



# Reference
https://www.designgurus.io/answers/detail/what-are-the-key-considerations-for-designing-a-microservices-architecture
