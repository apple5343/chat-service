FROM golang:1.25.5-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/main/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]