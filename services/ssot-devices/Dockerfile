# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy shared module
COPY shared /app/shared

# Copy service files
COPY services/ssot-devices/go.mod services/ssot-devices/go.sum ./services/ssot-devices/

# Create minimal go.work for this build
RUN printf 'go 1.22\n\nuse (\n\t./shared\n\t./services/ssot-devices\n)\n' > go.work

WORKDIR /app/services/ssot-devices
RUN go mod download

COPY services/ssot-devices .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/ssot_devices

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata wget

COPY --from=builder /bin/service /bin/service

EXPOSE 8082

ENTRYPOINT ["/bin/service"]
