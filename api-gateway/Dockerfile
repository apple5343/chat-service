FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/api_gateway/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]