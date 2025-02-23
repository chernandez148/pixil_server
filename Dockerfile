# Use an official Golang image as a build stage
FROM golang:1.20 AS builder

WORKDIR /app

# Copy go.mod and go.sum files for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o main .

# Use a lightweight image for running the app
FROM alpine:latest  

WORKDIR /root/

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your app runs on
EXPOSE 8080  

# Start the Go app
CMD ["./main"]
