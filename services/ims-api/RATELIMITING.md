# Rate Limiting Implementation (SE-002)

This document describes the rate limiting middleware implementation for the IMS API service.

## Overview

The rate limiting middleware protects the API from abuse by limiting the number of requests per tenant or IP address within a time window. It uses Redis for distributed rate limiting with a sliding window algorithm for accurate counting.

## Architecture

### Components

1. **Middleware** (`internal/middleware/ratelimit.go`)
   - `RateLimit()` - Rate limits by tenant ID (for authenticated requests)
   - `RateLimitByIP()` - Rate limits by client IP (for unauthenticated endpoints)
   - Sliding window counter algorithm using Redis
   - Fail-open strategy (allows requests on Redis errors)

2. **Configuration** (`internal/config/config.go`)
   - `RateLimitEnabled` - Enable/disable rate limiting
   - `RateLimitReadRPM` - Requests per minute for read operations (default: 300)
   - `RateLimitWriteRPM` - Requests per minute for write operations (default: 100)
   - `RateLimitBurst` - Additional burst capacity (default: 50)

3. **Server Integration** (`internal/api/server.go`)
   - Applied globally to all `/v1` routes with read limits
   - Applied to write operations with stricter limits
   - Health check endpoints (`/healthz`, `/readyz`) are exempt

## Configuration

### Environment Variables

```bash
# Enable/disable rate limiting
RATE_LIMIT_ENABLED=true

# Read operations limit (requests per minute)
RATE_LIMIT_READ_RPM=300

# Write operations limit (requests per minute)
RATE_LIMIT_WRITE_RPM=100

# Burst capacity (additional requests allowed)
RATE_LIMIT_BURST=50
```

### Example Configurations

#### Production (Strict)
```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_READ_RPM=200
RATE_LIMIT_WRITE_RPM=50
RATE_LIMIT_BURST=25
```

#### Development (Relaxed)
```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_READ_RPM=1000
RATE_LIMIT_WRITE_RPM=500
RATE_LIMIT_BURST=100
```

#### Testing (Disabled)
```bash
RATE_LIMIT_ENABLED=false
```

## How It Works

### Sliding Window Algorithm

The implementation uses a sliding window counter algorithm for accurate rate limiting:

1. Redis stores counters in 1-minute windows
2. Current request count = current window count + (previous window count × overlap percentage)
3. This smooths out rate limiting at window boundaries
4. Keys automatically expire after 2 minutes to save memory

### Redis Key Format

```
ratelimit:tenant:{tenantID}:{endpoint}:{timestamp}
ratelimit:ip:{clientIP}:{endpoint}:{timestamp}
```

### Endpoint Normalization

Endpoints are normalized to group similar requests:
- `/v1/incidents/abc-123` → `v1/incidents/{id}`
- `/v1/work-orders/456` → `v1/work-orders/{id}`

This prevents rate limit bypass by using different IDs.

## Response Headers

All responses include rate limit headers:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000000
```

- `X-RateLimit-Limit` - Maximum requests per minute
- `X-RateLimit-Remaining` - Requests remaining in current window
- `X-RateLimit-Reset` - Unix timestamp when the limit resets

## Rate Limit Exceeded Response

When rate limit is exceeded, the API returns HTTP 429:

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 42

{
  "error": "rate_limit_exceeded",
  "message": "Too many requests. Please try again later.",
  "retryAfter": 42
}
```

## Rate Limit Tiers

### Read Operations (300 RPM default)
- `GET /v1/incidents`
- `GET /v1/work-orders`
- `GET /v1/attachments`
- All other GET endpoints

### Write Operations (100 RPM default)
- `POST /v1/incidents`
- `POST /v1/work-orders`
- `PATCH /v1/incidents/{id}/status`
- `PATCH /v1/work-orders/{id}/status`
- All POST, PATCH, PUT, DELETE operations

### No Rate Limit
- `/healthz`
- `/readyz`

## Client Implementation

### Handling Rate Limits

Clients should:

1. **Check rate limit headers** on every response
2. **Back off when approaching limit** (use `X-RateLimit-Remaining`)
3. **Respect retry-after** on 429 responses
4. **Implement exponential backoff** for retries

### Example (JavaScript)

```javascript
async function makeRequest(url, options) {
  const response = await fetch(url, options);

  // Check rate limit headers
  const limit = response.headers.get('X-RateLimit-Limit');
  const remaining = response.headers.get('X-RateLimit-Remaining');

  if (remaining < 10) {
    console.warn(`Rate limit warning: ${remaining}/${limit} remaining`);
  }

  if (response.status === 429) {
    const retryAfter = parseInt(response.headers.get('Retry-After') || '60');
    console.error(`Rate limited. Retry after ${retryAfter} seconds`);
    await sleep(retryAfter * 1000);
    return makeRequest(url, options); // Retry
  }

  return response;
}
```

### Example (Go)

```go
func makeRequest(client *http.Client, req *http.Request) (*http.Response, error) {
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }

    // Check rate limit headers
    remaining := resp.Header.Get("X-RateLimit-Remaining")
    if r, _ := strconv.Atoi(remaining); r < 10 {
        log.Printf("Rate limit warning: %s requests remaining", remaining)
    }

    if resp.StatusCode == http.StatusTooManyRequests {
        retryAfter := resp.Header.Get("Retry-After")
        seconds, _ := strconv.Atoi(retryAfter)
        time.Sleep(time.Duration(seconds) * time.Second)
        return makeRequest(client, req) // Retry
    }

    return resp, nil
}
```

## Monitoring

### Metrics to Monitor

1. **Rate limit hits** - How often 429 responses are returned
2. **Rate limit by tenant** - Which tenants are hitting limits
3. **Redis performance** - Latency and error rates
4. **Rate limit bypass attempts** - Suspicious patterns

### Logging

The middleware logs:
- Rate limit exceeded events (WARN level)
- Redis errors (ERROR level)
- Missing tenant/IP context (WARN level)

Example log:
```json
{
  "level": "warn",
  "msg": "rate limit exceeded",
  "tenant_id": "demo-tenant",
  "endpoint": "v1/incidents/{id}",
  "retry_after": 45
}
```

## Testing

### Manual Testing

1. **Test read limits:**
```bash
# Make 301 requests in 1 minute
for i in {1..301}; do
  curl -H "X-Tenant-Id: demo-tenant" http://localhost:8080/v1/incidents
done
```

2. **Test write limits:**
```bash
# Make 101 POST requests in 1 minute
for i in {1..101}; do
  curl -X POST -H "X-Tenant-Id: demo-tenant" \
    -H "Content-Type: application/json" \
    -d '{"title":"test"}' \
    http://localhost:8080/v1/incidents
done
```

3. **Verify headers:**
```bash
curl -v http://localhost:8080/v1/incidents | grep X-RateLimit
```

### Load Testing

Use tools like `wrk` or `k6` for load testing:

```bash
# Test with wrk
wrk -t4 -c100 -d60s \
  -H "X-Tenant-Id: demo-tenant" \
  http://localhost:8080/v1/incidents
```

## Troubleshooting

### Rate limit not working

1. Check `RATE_LIMIT_ENABLED=true` is set
2. Verify Redis is running and accessible
3. Check tenant ID is present in context
4. Review logs for Redis errors

### Too many 429 errors

1. Increase limits via environment variables
2. Check if burst capacity is sufficient
3. Review client retry logic
4. Consider per-tenant overrides (future enhancement)

### Redis errors

1. Check Redis connectivity
2. Verify Redis memory limits
3. Monitor Redis performance
4. Consider Redis cluster for scale

## Future Enhancements

1. **Per-tenant rate limit overrides** - Custom limits for specific tenants
2. **Dynamic rate limits** - Adjust limits based on system load
3. **Rate limit analytics** - Dashboard for monitoring usage
4. **Whitelist/blacklist** - Bypass or block specific tenants/IPs
5. **Multi-tier limits** - Different limits by subscription tier
6. **Distributed tracing** - Integration with OpenTelemetry

## References

- Implementation: `/home/pato/opt/ESSP/services/ims-api/internal/middleware/ratelimit.go`
- Configuration: `/home/pato/opt/ESSP/services/ims-api/internal/config/config.go`
- Server setup: `/home/pato/opt/ESSP/services/ims-api/internal/api/server.go`
- Redis algorithms: [Sliding Window Rate Limiting](https://redis.io/glossary/rate-limiting/)
