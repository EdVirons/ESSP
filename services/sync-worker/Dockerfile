# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy shared module
COPY shared /app/shared

# Copy service files
COPY services/sync-worker/go.mod services/sync-worker/go.sum ./services/sync-worker/

# Create minimal go.work for this build
RUN printf 'go 1.22\n\nuse (\n\t./shared\n\t./services/sync-worker\n)\n' > go.work

WORKDIR /app/services/sync-worker
RUN go mod download

COPY services/sync-worker .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/worker

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata wget

COPY --from=builder /bin/service /bin/service

EXPOSE 8084

ENTRYPOINT ["/bin/service"]
