# Stage 1: Build the application
FROM golang:1.22 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better cache utilization
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 2: Create a minimal image
FROM gcr.io/distroless/static:nonroot

# Set working directory
WORKDIR /

# Copy binary from builder
COPY --from=builder /app/server .

# Run the binary
ENTRYPOINT ["/server"]
