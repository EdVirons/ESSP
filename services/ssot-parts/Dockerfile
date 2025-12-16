# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy go workspace files from root
COPY ../../go.work ../../go.work.sum ./

# Copy shared module
COPY ../../shared /app/shared

# Copy service files
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/ssot_parts

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /bin/service /bin/service

EXPOSE 8080

ENTRYPOINT ["/bin/service"]
