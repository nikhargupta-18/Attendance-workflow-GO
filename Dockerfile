FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o attendance-api ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/attendance-api .

EXPOSE 8080

CMD ["./attendance-api"]
