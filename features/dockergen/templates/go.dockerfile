# Build stage: Compile the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache dependencies first for faster rebuilds
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Final stage: Create a minimal runtime image
FROM alpine:3.21 

WORKDIR /root/

# Copy the statically linked binary from the builder stage
COPY --from=builder /app/main .

# Optionally, add a non-root user for security (recommended for production)
# RUN adduser -D appuser
# USER appuser

CMD ["./main"]
