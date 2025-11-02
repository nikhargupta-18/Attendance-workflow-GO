FROM golang:1.21-alpine AS builder

WORKDIR /app

# Add git and build dependencies
RUN apk add --no-cache git make build-base

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o attendance-api ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/attendance-api .
COPY .env .

EXPOSE 8080

CMD ["./attendance-api"]