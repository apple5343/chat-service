FROM golang:alpine

WORKDIR /app

COPY . .

CMD ["sh", "-c", "go test -v ./test/integration/... "]