# Build dashboard stage
FROM node:22-alpine AS dashboard-builder
WORKDIR /dashboard

# Copy dashboard source
COPY dashboard/package.json dashboard/package-lock.json* ./
RUN npm ci

COPY dashboard ./
RUN npm run build

# Build Go service stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy shared module
COPY shared /app/shared

# Copy service files
COPY services/ims-api/go.mod services/ims-api/go.sum ./services/ims-api/

# Create minimal go.work for this build (only shared + ims-api)
RUN printf 'go 1.22\n\nuse (\n\t./shared\n\t./services/ims-api\n)\n' > go.work

WORKDIR /app/services/ims-api
RUN go mod download

COPY services/ims-api .

# Copy dashboard build output to embed directory
COPY --from=dashboard-builder /dashboard/dist ./internal/admin/dashboard/dist/

# Build the service
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/service ./cmd/api

# Runtime stage
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /bin/service /bin/service

EXPOSE 8080

ENTRYPOINT ["/bin/service"]
