# Start from the official Go image for the build stage
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire repository into the build context
COPY . .

# Set the working directory to the CustomerService
WORKDIR /app/services/customerService

# Replace the broker module with the local path
RUN go mod edit -replace github.com/rasm445f/soft-exam-2/broker=../../broker

# Download dependencies
RUN go mod download

# Build the Go application
RUN go build -o /api

# Expose the service port
EXPOSE 8081

# Set the default command to run the application
CMD ["/api"]
