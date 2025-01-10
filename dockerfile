# Use the official Go image as the base image
FROM golang:1.23.4-alpine3.19 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o myiris-app .

# Use a lightweight alpine image for the final stage
FROM alpine:3.19

# Install necessary certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/myiris-app .

# Copy .env file if it exists (using shell trick to make it optional)
COPY .env* ./

# Expose the port the app runs on (adjust as needed)
EXPOSE 8080

# Command to run the executable
CMD ["./myiris-app"]