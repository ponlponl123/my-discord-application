FROM golang:alpine AS builder

WORKDIR /app
COPY . .

# Build with options to create a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o main src/main.go

# Use a scratch image (smallest possible base)
FROM scratch

# Copy only the binary from builder
COPY --from=builder /app/main /main

# Command to run
CMD ["/main"]