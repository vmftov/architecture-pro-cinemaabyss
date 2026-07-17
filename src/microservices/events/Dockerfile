# Этап сборки
FROM golang:1.26.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o events-service .

# Use a smaller image for the final application
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/events-service /app/events-service

WORKDIR /app

# Expose the application port
EXPOSE 8082

# Command to run the executable
CMD ["./events-service"]