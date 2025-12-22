# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy shared module
COPY shared /app/shared

# Copy service files
COPY services/ssot-hr/go.mod services/ssot-hr/go.sum ./services/ssot-hr/

# Create minimal go.work for this build
RUN printf 'go 1.22\n\nuse (\n\t./shared\n\t./services/ssot-hr\n)\n' > go.work

WORKDIR /app/services/ssot-hr
RUN go mod download

COPY services/ssot-hr .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/ssot_hr

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata wget

COPY --from=builder /bin/service /bin/service

EXPOSE 8300

ENTRYPOINT ["/bin/service"]
