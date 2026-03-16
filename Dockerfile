FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o scheduler ./cmd/main.go

FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates libsqlite3-0 && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/scheduler .
COPY web/ ./web/
ENV TODO_PORT=7540
EXPOSE 7540
CMD ["./scheduler", "web/"]
