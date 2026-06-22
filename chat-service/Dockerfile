FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/main/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]