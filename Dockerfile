FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN go mod download

COPY . .

RUN go build -o main ./main.go

FROM alpine:latest AS production

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main.go"]