# Build
FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o auth-service ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-service .

RUN mkdir -p /shared/uploads/avatar

EXPOSE 8080

CMD ["./auth-service"]