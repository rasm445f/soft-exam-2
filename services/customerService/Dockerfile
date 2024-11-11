# Start from the official Go image for the build stage
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./


# Build the Go application
RUN go build -o /api
# RUN CGO_ENABLED=0 GOOS=linux go build -o /hello

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# # Start a new stage from a lightweight base image
# FROM alpine:latest
#
# # Copy the compiled binary from the build stage
# COPY --from=builder /app/hello /hello

# Run the hello binary
CMD ["/api"]
