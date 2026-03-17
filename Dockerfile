FROM golang:1.24-alpine3.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/main.go

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/scheduler .
COPY web/ ./web/
CMD ["./scheduler", "web/"]
