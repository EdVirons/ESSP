# Rate Limiting Quick Start Guide

## Overview

Rate limiting is enabled by default for the IMS API service to protect against abuse and ensure fair resource usage.

## Default Configuration

```bash
RATE_LIMIT_ENABLED=true
RATE_LIMIT_READ_RPM=300      # 300 requests per minute for GET operations
RATE_LIMIT_WRITE_RPM=100     # 100 requests per minute for POST/PATCH/PUT/DELETE
RATE_LIMIT_BURST=50          # Additional burst capacity
```

## How It Works

### Two-Tier System

1. **Read Operations (300 RPM)** - Applied to all GET requests
2. **Write Operations (100 RPM)** - Applied to all POST, PATCH, PUT, DELETE requests

### Per-Tenant Rate Limiting

- Each tenant has their own rate limit bucket
- Limits are tracked independently per tenant
- Uses tenant ID from the `X-Tenant-Id` header

### Response Headers

Every response includes:
```
X-RateLimit-Limit: 300           # Maximum requests allowed
X-RateLimit-Remaining: 250       # Requests remaining in current window
X-RateLimit-Reset: 1640000000    # Unix timestamp when limit resets
```

## Rate Limit Exceeded

When you exceed the limit:

**HTTP Response:**
```http
HTTP/1.1 429 Too Many Requests
Retry-After: 42
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640000042

{
  "error": "rate_limit_exceeded",
  "message": "Too many requests. Please try again later.",
  "retryAfter": 42
}
```

## Client Implementation

### Check Headers Before Making Requests

```javascript
// Check remaining requests before proceeding
const remaining = parseInt(lastResponse.headers['x-ratelimit-remaining']);
const limit = parseInt(lastResponse.headers['x-ratelimit-limit']);

if (remaining < limit * 0.1) {  // Less than 10% remaining
  console.warn('Approaching rate limit, slowing down requests');
  await sleep(1000);
}
```

### Handle 429 Responses

```javascript
if (response.status === 429) {
  const retryAfter = parseInt(response.headers['retry-after']) || 60;
  console.log(`Rate limited. Waiting ${retryAfter} seconds...`);
  await sleep(retryAfter * 1000);
  // Retry the request
}
```

## Common Scenarios

### Scenario 1: Bulk Data Import

**Problem:** Need to import 1000 records

**Solution:**
- Batch requests: Send 10 records per request = 100 requests
- Rate limit: 100 RPM for writes
- Time needed: ~1 minute
- Implementation: Add 600ms delay between requests

```javascript
for (let batch of batches) {
  await postBatch(batch);
  await sleep(600);  // 600ms between requests = 100 RPM
}
```

### Scenario 2: Real-Time Monitoring Dashboard

**Problem:** Dashboard polls every 5 seconds for updates

**Solution:**
- Read limit: 300 RPM
- 5-second polling = 12 requests per minute
- Safe within limits (4% of capacity)

### Scenario 3: Multiple Services Sharing Tenant

**Problem:** Multiple services using same tenant ID

**Solution:**
- Coordinate across services
- Implement request queuing
- Consider increasing limits via config
- Monitor X-RateLimit-Remaining header

## Adjusting Limits

### Increase Limits (e.g., for development)

```bash
export RATE_LIMIT_READ_RPM=1000
export RATE_LIMIT_WRITE_RPM=500
export RATE_LIMIT_BURST=100
```

### Disable Rate Limiting (e.g., for local testing)

```bash
export RATE_LIMIT_ENABLED=false
```

### Production Tuning

1. Monitor 429 error rate
2. Check p95/p99 latency
3. Review tenant usage patterns
4. Adjust limits based on data

```bash
# Conservative (high security)
RATE_LIMIT_READ_RPM=100
RATE_LIMIT_WRITE_RPM=50
RATE_LIMIT_BURST=25

# Moderate (default)
RATE_LIMIT_READ_RPM=300
RATE_LIMIT_WRITE_RPM=100
RATE_LIMIT_BURST=50

# Relaxed (high throughput)
RATE_LIMIT_READ_RPM=600
RATE_LIMIT_WRITE_RPM=200
RATE_LIMIT_BURST=100
```

## Exempt Endpoints

These endpoints are NOT rate limited:
- `/healthz` - Health check
- `/readyz` - Readiness check

## Troubleshooting

### Getting 429 errors unexpectedly?

1. **Check your request rate:**
   ```bash
   # Count requests in last minute
   grep "your-tenant-id" /var/log/ims-api.log | \
     grep $(date -u +%Y-%m-%dT%H:%M) | wc -l
   ```

2. **Check rate limit headers:**
   ```bash
   curl -v http://api/v1/incidents \
     -H "X-Tenant-Id: your-tenant" 2>&1 | grep X-RateLimit
   ```

3. **Verify configuration:**
   ```bash
   # Check what limits are configured
   env | grep RATE_LIMIT
   ```

### Not seeing rate limit headers?

- Check that rate limiting is enabled
- Verify Redis is running
- Check tenant ID is being set correctly
- Review application logs for errors

## Best Practices

1. **Always check rate limit headers** on every response
2. **Implement exponential backoff** when hitting 429s
3. **Monitor remaining requests** and slow down proactively
4. **Batch operations** when possible to reduce request count
5. **Cache read responses** to minimize redundant API calls
6. **Use webhooks/events** instead of polling when available
7. **Implement request queuing** for high-volume operations

## Monitoring Queries

### Check rate limit hits by tenant

```bash
# Count 429 responses by tenant
grep "rate limit exceeded" /var/log/ims-api.log | \
  jq -r '.tenant_id' | sort | uniq -c | sort -rn
```

### Calculate request rate

```bash
# Requests per minute for a tenant
grep "demo-tenant" /var/log/ims-api.log | \
  grep $(date -u +%Y-%m-%dT%H:%M) | wc -l
```

## Support

For questions or issues with rate limiting:

1. Check this guide first
2. Review full documentation: `RATELIMITING.md`
3. Check logs: `grep "rate limit" /var/log/ims-api.log`
4. Contact platform team with specific examples

## Technical Details

- **Algorithm:** Sliding window counter
- **Storage:** Redis (keys expire after 2 minutes)
- **Granularity:** Per-minute windows
- **Accuracy:** ~95% (due to sliding window approximation)
- **Performance Impact:** ~1-2ms per request
- **Fail Mode:** Fail-open (allows requests if Redis is down)

## Changes from Previous Version

If upgrading from a version without rate limiting:

1. Rate limiting is **enabled by default**
2. Existing clients will see new response headers
3. High-volume clients may need to adjust request patterns
4. Health check endpoints remain unlimited
5. Per-tenant isolation prevents one tenant affecting others

## Example: Compliant Client

```python
import requests
import time

class RateLimitedClient:
    def __init__(self, base_url, tenant_id):
        self.base_url = base_url
        self.headers = {'X-Tenant-Id': tenant_id}
        self.remaining = None
        self.limit = None
        self.reset_time = None

    def request(self, method, path, **kwargs):
        # Check if we're close to limit
        if self.remaining is not None and self.remaining < 10:
            wait_time = max(0, self.reset_time - time.time())
            if wait_time > 0:
                print(f"Approaching limit, waiting {wait_time:.1f}s")
                time.sleep(wait_time)

        # Make request
        url = f"{self.base_url}{path}"
        response = requests.request(method, url, headers=self.headers, **kwargs)

        # Update rate limit state
        self.limit = int(response.headers.get('X-RateLimit-Limit', 0))
        self.remaining = int(response.headers.get('X-RateLimit-Remaining', 0))
        self.reset_time = int(response.headers.get('X-RateLimit-Reset', 0))

        # Handle 429
        if response.status_code == 429:
            retry_after = int(response.headers.get('Retry-After', 60))
            print(f"Rate limited. Retrying after {retry_after}s")
            time.sleep(retry_after)
            return self.request(method, path, **kwargs)  # Retry

        return response

    def get(self, path, **kwargs):
        return self.request('GET', path, **kwargs)

    def post(self, path, **kwargs):
        return self.request('POST', path, **kwargs)

# Usage
client = RateLimitedClient('http://localhost:8080', 'demo-tenant')
response = client.get('/v1/incidents')
print(f"Remaining requests: {client.remaining}/{client.limit}")
```
