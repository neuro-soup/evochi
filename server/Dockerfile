FROM golang:1.23-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /server ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /server .

EXPOSE 8080

CMD ["./server"]
