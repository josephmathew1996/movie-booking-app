# Stage 1: Build the Go binary
# official Go base image that has all tools and packages to compile and run a Go application 
FROM docker-public.docker.devstack.vwgroup.com/golang:1.22-alpine AS builder

# create a directory inside the image.
# instructs Docker to use this directory as the default destination for all subsequent commands
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/userservice ./cmd/api

# Stage 2: Create a small image with the Go binary
FROM docker-public.docker.devstack.vwgroup.com/alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /app/userservice .

# exposes port `8080` for accessing the api server
EXPOSE 8080

# tell Docker what command to execute when our image is used to start a container
CMD [ "./userservice" ]