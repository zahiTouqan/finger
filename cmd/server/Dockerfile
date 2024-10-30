FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY cmd/server/main.go ./cmd/server/
COPY finger.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /server .
EXPOSE 79
CMD ["./server"]