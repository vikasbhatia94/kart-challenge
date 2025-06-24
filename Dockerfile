# Stage 1: Build the application
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker's layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code.
COPY . .

# Build the Go application as a static binary.
# This enables it to run in a minimal container without C libraries.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /kart ./main.go

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Create a non-root user and group for security.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create app directory and set ownership.
RUN mkdir /app && chown -R appuser:appgroup /app

# Switch to the non-root user.
USER appuser
WORKDIR /app

# Copy just the built binary from the builder stage.
COPY --from=builder /kart .

# Expose the port the server runs on.
EXPOSE 8080

# The command to run the application.
CMD ["./kart"] 