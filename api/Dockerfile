# First stage: build the Go application
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install swag for Swagger doc generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code
COPY . .

# Generate Swagger docs
RUN swag init -g cmd/server/main.go -o docs

# Build the Go application
RUN go build -o /app/bin/app ./cmd/server

# Second stage: use a minimal base image
FROM alpine:3.20

# Install curl for healthcheck
RUN apk add --no-cache curl

# Copy the built Go binary from the builder stage
COPY --from=builder /app/bin/app /app/bin/app
COPY --from=builder /app/docs /app/docs

ENTRYPOINT ["/app/bin/app"]
