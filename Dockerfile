# Use the official Go image to build the application
FROM golang:1.23.3-bookworm AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY main.go go.mod go.sum ./
RUN go mod download

# Build the Go application
RUN go build -o app .

# Use a smaller base image for the final image
FROM debian:bookworm-slim

# Copy the binary from the builder stage
COPY --from=builder /app/app /app/app

# Set the working directory
WORKDIR /app

EXPOSE 3000

# Command to run the application
CMD ["./app"]
