FROM golang:1.23.2-alpine as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY cmd/client/main.go ./cmd/client/
RUN CGO_ENABLED=0 GOOS=linux go build -o /client ./cmd/client

FROM alpine:latest
WORKDIR /app
COPY --from=builder /client .
ENTRYPOINT ["./client"]