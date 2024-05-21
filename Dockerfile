# Use the official Golang image version 1.20.1 for building the application
FROM golang:1.21.10-alpine3.18 AS builder

# Create a directory for the application
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main

# Expose the port that the Gin webserver will run on
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
