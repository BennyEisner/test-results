# First stage: build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o /app/bin/app ./cmd/server

# Second stage: use a minimal base image
FROM alpine:3.20

# Install curl for healthcheck
RUN apk add --no-cache curl

# Copy the built Go binary from the builder stage
COPY --from=builder /app/bin/app /app/bin/app

ENTRYPOINT ["/app/bin/app"]
