# User service for Movie Booking App
This application manages the user service.

## Project structure 
Project structure is inspired by the the `Standard Go Project Layout` explained in the below repository.

Reference: https://github.com/golang-standards/project-layout

## Prerequisites

1. Install Golang SDK: https://go.dev/dl/
2. Execute following commands:
     
### To install all the dependencies

        go mod tidy -v 
    
### To run the server
    
        go run cmd/api/main.go

## Libraries used

1. Zap for logging: https://github.com/uber-go/zap
2. Viper for managing environment variables: https://github.com/spf13/viper
3. Echo framework for setting up the server: https://echo.labstack.com/
4. Standard libraries for other basic purposes: https://pkg.go.dev/std

## Test server

### Sample Commands


## setup docker

docker build -t user-service .

docker run -p 8080:8080 user-service
