FROM golang:alpine AS builder

WORKDIR /app
COPY . .

# Install CA certificates in builder stage
RUN apk add --no-cache ca-certificates

# Build with options to create a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o main src/main.go

# Use a scratch image (smallest possible base)
FROM scratch

# Copy the CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy only the binary from builder
COPY --from=builder /app/main /main

# Command to run
CMD ["/main"]
