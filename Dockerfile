FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ ./
COPY *.yaml .

EXPOSE 8080

CMD ["./server"]