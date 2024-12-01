# CNAD_Assignment1


# Design consideration of your microservices

Service Decomposition:
Login : retrive user id and password to login
Signup : Create user account
Profile: Manage user profile , view past booking and manage existing booking
Car Listing: View all listing of car
Reservation: Create reservation
Checkout: caculate amount and create receipt

Inter-Service communication
User HTTP/REST for synchronous call (eg from car listing to reservation)

API Gateway
Acts as the entry point for all client requests and helps route requests to the appropriate microservice.

Security
Secure all sensetive data such as user password using JWT Token


# Architecture diagram
![architecture diagram](https://github.com/user-attachments/assets/57d6a3d3-3063-4e9e-9ace-a66fd00122d4)

# Instructions for setting up and running your microservices





# Reference
https://www.designgurus.io/answers/detail/what-are-the-key-considerations-for-designing-a-microservices-architecture
