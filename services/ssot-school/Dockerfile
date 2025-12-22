# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy shared module
COPY shared /app/shared

# Copy service files
COPY services/ssot-school/go.mod services/ssot-school/go.sum ./services/ssot-school/

# Create minimal go.work for this build
RUN printf 'go 1.22\n\nuse (\n\t./shared\n\t./services/ssot-school\n)\n' > go.work

WORKDIR /app/services/ssot-school
RUN go mod download

COPY services/ssot-school .

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/ssot_school

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata wget

COPY --from=builder /bin/service /bin/service

EXPOSE 8081

ENTRYPOINT ["/bin/service"]
