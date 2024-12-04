# Start from the official Go image for the build stage
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files from the customerService directory
COPY services/customerService/go.mod services/customerService/go.sum ./

# Copy the broker directory again to ensure it is available for go mod download
COPY broker ./broker

# RUN go mod download
RUN go mod edit -replace github.com/rasm445f/soft-exam-2/broker=./broker
RUN go mod download

# Copy the rest of the customerService directory
COPY services/customerService/ ./

# Build the Go application
RUN go build -o /api
# RUN CGO_ENABLED=0 GOOS=linux go build -o /hello

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8081

# # Start a new stage from a lightweight base image
# FROM alpine:latest
#
# # Copy the compiled binary from the build stage
# COPY --from=builder /app/hello /hello

# Run the hello binary
CMD ["/api"]
