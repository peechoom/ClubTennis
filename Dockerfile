# Use the official Golang image version 1.21.10 for building the application
FROM golang:1.21.10-alpine3.18 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create a directory for the application
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the entire source code to the container
COPY . .

# Build the Go application
RUN go build -o main

# Expose the port defined by the PORT environment variable
# Cloud Run provides the PORT environment variable, so we use that
EXPOSE $PORT

# Command to run the binary
CMD ["./main"]
