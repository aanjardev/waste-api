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

# Install postgresql-client untuk psql & pg_isready
RUN apk add --no-cache ca-certificates tzdata postgresql16-client

COPY --from=builder /app/bin/waste-api /app/waste-api
COPY --from=builder /app/migrations /app/migrations

RUN mkdir -p /app/uploads

# Set executable permission untuk script
RUN chmod +x /app/migrations/apply.sh

EXPOSE 8080

ENTRYPOINT ["/app/migrations/apply.sh"]
CMD ["/app/waste-api"]