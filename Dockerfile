# ---------- Build stage ----------
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Install git & ca-certificates untuk dependency
RUN apk add --no-cache git ca-certificates

# Copy go.mod & go.sum
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/waste-api ./cmd/api

# ---------- Runtime stage ----------
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Copy binary dari builder
COPY --from=builder /app/bin/waste-api /app/waste-api

# Copy folder migrations
COPY --from=builder /app/migrations /app/migrations

# Buat folder uploads
RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["/app/waste-api"]