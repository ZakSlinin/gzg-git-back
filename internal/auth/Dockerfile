# Build
FROM golang:alpine:latest AS builder

# Install Git
RUN apk add --no-cache git

# Work directory
WORKDIR /app

# Copy go mod and sum
COPY go.mod go.sum /app/

# Download dependencies
RUN go mod download

# Copy auth_service
COPY internal/auth /app/internal/auth

# Build binary
RUN go build -o auth-service /app/internal/auth/main.go

FROM alpine:latest

WORKDIR /app

# Copy from builder
COPY --from=builder /app/auth-service .

# Create directory to avatars
RUN mkdir -p uploads/avatar

EXPOSE 8081

CMD ["./auth-service"]
