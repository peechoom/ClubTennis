# Use the official Golang image version 1.21.10 for building the application
FROM golang:1.21.10-alpine3.18 AS builder

# Set environment variables
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Create a directory for the application
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Download libvps dependancy for webp compression (required by https://github.com/h2non/bimg)
RUN apk add --no-cache vips vips-dev gcc musl-dev

# Copy the entire source code to the container
COPY . .

# ensure that scripts are executable
RUN chmod -x ./scripts/wait_for_it.sh

# Build the Go application
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main



# Define the final image for running the application
FROM alpine:3.18
RUN apk add --no-cache bash
RUN apk add --no-cache vips vips-dev
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config/.env ./config/.env
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/scripts/wait_for_it.sh ./scripts/wait_for_it.sh
RUN chmod +x ./scripts/wait_for_it.sh

HEALTHCHECK --interval=5s --timeout=2s --start-period=3s --retries=3 CMD [ "nc", "-z", "server", "8080" ]