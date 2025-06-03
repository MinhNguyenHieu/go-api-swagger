FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/docs ./docs

COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./main"]

