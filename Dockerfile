FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o cmd/main ./cmd/main.go

FROM alpine:latest AS production

WORKDIR /app

COPY --from=builder /app/cmd/main .

EXPOSE 8080

CMD ["./cmd/main.go"]