# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o app ./main.go

# Final stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./app", "serve"]
