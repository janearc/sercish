# Use an official Golang image as the base image for building the application
FROM golang:1.18-alpine AS builder

ENV GOOS=linux
ENV GOARCH=amd64

ENV APP_ROOT=/app

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# start building
RUN go build -o bitey .

# Use a lightweight image to run the application
FROM alpine:latest

ENV APP_ROOT=/app

# Copy the built Go binary from the builder image
COPY --from=builder /app/archeavy .
COPY --from=builder /app/public /app/public
COPY --from=builder /app/views /app/views

# Expose the application's port
EXPOSE 8080

# Command to run the application
CMD ["./bitey"]
