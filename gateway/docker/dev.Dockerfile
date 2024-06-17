# Use the official Golang image as the base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["go","run","main.go"]