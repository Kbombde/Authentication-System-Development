# Authentication-System-Development
This repository contains the backend code for an authentication system developed in Golang using the Gin framework and PostgreSQL as the database. The system allows user sign-up, login, password reset, and email verification functionalities.
Features
User Signup with email and password.
User Login with email and password authentication.
Password Reset functionality.
Email Verification to validate user email addresses.
API-based architecture for ease of integration.

Steps to Set Up Environment
Install Golang:
Download and install Golang from golang.org.
Set Up PostgreSQL Database:
Install PostgreSQL.
Create a database: CREATE DATABASE auth_system;.
Run the Server:
go run main.go
Testing:
Use the provided Postman collection for API testing.

API Endpoints
Base URL
http://localhost:8080

Endpoint	                 Method	                 Description
/signup	                   POST	               Register user with email and password
/login	                   POST	               Log in with email and password
/update-password	         PUT	               Update the password for a user
/reset-password           POST	               Reset password for a user
/verify-email           	POST	               Mark an email as verified
/verify-password	        POST	               Verify email and password combination
