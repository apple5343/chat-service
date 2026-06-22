FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/project_service/main.go

# Stage 2: Prepare the final image
FROM alpine:latest

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main /main

# Set the entrypoint command
CMD ["/main"]
