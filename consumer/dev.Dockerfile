# Use the official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and install the Go dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
# RUN go build -o healthcheck

# Expose the port that the application listens on
# EXPOSE 8080

# Run the Go application by go run command
CMD ["go", "run", "main.go", "consumer.go", "server.go", "database.go" ,"kafka.go"]