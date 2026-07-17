# Этап сборки
FROM golang:1.26.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o proxy .

# Use a smaller image for the final application
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/proxy /app/proxy

WORKDIR /app

# Expose the application port
EXPOSE 8000

# Command to run the executable
CMD ["./proxy"]