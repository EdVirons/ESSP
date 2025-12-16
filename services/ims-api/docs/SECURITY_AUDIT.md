# Security Audit Report - IMS API Service (SE-001)

**Audit Date:** 2025-12-12
**Service:** ims-api
**Repository:** /home/pato/opt/ESSP/services/ims-api
**Auditor:** Security Team

---

## Executive Summary

This security audit was performed on the IMS (Incident Management System) API service to identify and remediate potential security vulnerabilities. The audit covered SQL injection risks, authentication/authorization controls, input validation, secret handling, and general security best practices.

**Overall Security Posture:** GOOD with CRITICAL FIX APPLIED
**Vulnerabilities Found:** 1 Critical Bug, 0 High Risk, 2 Medium Risk
**Remediation Status:** All issues addressed

---

## 1. SQL Injection Review

### Findings

**Status:** PASS

All SQL queries in the `internal/store/` directory use parameterized queries with PostgreSQL's pgx driver. No instances of string concatenation or `fmt.Sprintf` were found in SQL construction.

**Files Audited:**
- `/home/pato/opt/ESSP/services/ims-api/internal/store/workorders_repo.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/store/incidents_repo.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/store/workorder_parts_repo.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/store/schools_repo.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/store/parts_repo.go`
- `/home/pato/opt/ESSP/services/ims-api/internal/store/inventory_repo.go`
- All other repository files in `internal/store/*_repo.go`

**Good Practices Observed:**
- Consistent use of positional parameters (`$1`, `$2`, etc.)
- Dynamic WHERE clause construction uses parameterized placeholders
- ILIKE queries properly parameterize the search term with wildcards added in Go code
- Example from `incidents_repo.go`:
  ```go
  if p.Query != "" {
      conds = append(conds, "(title ILIKE $"+itoa(argN)+" OR description ILIKE $"+itoa(argN)+")")
      args = append(args, "%"+p.Query+"%")
      argN++
  }
  ```

**Issues Found:** None

---

## 2. Authentication/Authorization Review

### Findings

**Status:** PASS with OBSERVATIONS

#### JWT Verification (`internal/auth/verifier.go`)

**Strengths:**
- Proper JWT signature verification using RS256 algorithm
- JWKS key rotation support with 10-minute cache
- Issuer (`iss`) validation
- Audience (`aud`) validation
- Algorithm restriction to RS256 only (prevents "none" algorithm attack)
- Concurrent-safe key caching with RWMutex

**Code Review:**
```go
parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
```

**Observations:**
- Token expiration validation is handled by the jwt library's default behavior
- No explicit check for `nbf` (not before) claim, but jwt library handles this
- Good error handling without exposing internal details

#### Authentication Middleware (`internal/middleware/auth.go`)

**Strengths:**
- Properly extracts and validates Bearer token format
- Sets tenant and school IDs from JWT claims to request headers
- Returns 401 for missing or invalid tokens

**Observation:**
- JWT claims `tenantId` and `schoolId` are trusted without additional validation
- This is acceptable as they come from verified tokens, but ensure the token issuer is properly validating these claims

#### Tenant Isolation (`internal/middleware/tenancy.go`)

**Strengths:**
- All queries include `tenant_id` and `school_id` filters
- Context-based tenant/school isolation
- Falls back to dev tenant/school when headers are missing (dev mode only)

**Example from `workorders_repo.go`:**
```go
WHERE tenant_id=$1 AND school_id=$2 AND id=$3
```

**Verified in Multiple Files:**
- `GetByID()` methods consistently check tenant + school + id
- `List()` methods filter by tenant and school
- `Update` operations verify ownership before modification

**Issues Found:** None - tenant isolation is properly implemented throughout

---

## 3. Input Validation Review

### Findings

**Status:** GOOD with RECOMMENDATIONS

#### Handler Validation (`internal/handlers/*.go`)

**Strengths:**
- Required field validation on all create endpoints
- Proper JSON decoding error handling
- Input trimming with `strings.TrimSpace()`
- Limit validation with max bounds (50 default, 200 max)

**Examples:**
```go
// From incidents.go
if strings.TrimSpace(req.DeviceID) == "" || strings.TrimSpace(req.Title) == "" {
    http.Error(w, "deviceId and title are required", http.StatusBadRequest)
    return
}

// From utils.go - limit parsing
func parseLimit(raw string, def, max int) int {
    // Safe integer parsing without strconv
    // Prevents negative values and enforces maximum
}
```

**Observations:**
- Custom cursor decoding with base64 validation
- No path traversal vulnerabilities in file operations (object keys are generated)
- URL parameters extracted safely via chi.URLParam

**Request Body Size:**
- Previously: No global limit (VULNERABILITY - Medium Risk)
- Fixed: 10MB limit now enforced via `MaxBodySize` middleware

**Recommendations:**
1. Consider adding field length validation for strings (e.g., title max 500 chars)
2. Add email format validation for contact emails
3. Consider phone number format validation

---

## 4. Secret Handling Review

### Findings

**Status:** PASS

#### Configuration (`internal/config/config.go`)

**Strengths:**
- All secrets loaded from environment variables
- Default values only for development/non-sensitive configs
- No hardcoded production credentials

**Sensitive Fields Properly Handled:**
- `PGDSN` - Database connection string
- `RedisPassword` - Redis password
- `MinIOAccessKey` - Object storage credentials
- `MinIOSecretKey` - Object storage credentials
- `AuthJWKSURL` - External auth endpoint

**Example:**
```go
MinIOAccessKey: getenv("MINIO_ACCESS_KEY", "minioadmin"),
MinIOSecretKey: getenv("MINIO_SECRET_KEY", "minioadmin"),
```

**Note:** Default "minioadmin" credentials are for local development only.

#### Logging (`internal/middleware/http.go`, `internal/handlers/*.go`)

**Strengths:**
- Error messages do not expose sensitive data
- Generic error messages returned to clients
- No credential logging observed
- Stack traces only logged on panic, not sent to client

**Examples:**
```go
http.Error(w, "failed to create incident", http.StatusInternalServerError)
// vs internal logging:
log.Error("panic recovered", zap.Any("panic", rec), zap.ByteString("stack", debug.Stack()))
```

**Issues Found:** None

---

## 5. Additional Security Findings

### Critical Bug Fixed

**File:** `/home/pato/opt/ESSP/services/ims-api/internal/store/incidents_repo.go`
**Severity:** CRITICAL
**Status:** FIXED

**Issue:**
The `Create()` method had a parameter mismatch between the INSERT column list (28 columns) and the values provided (14 values). This would cause all incident creation to fail.

**Before:**
```go
`, inc.ID, inc.TenantID, inc.SchoolID, inc.DeviceID, inc.Category, inc.Severity, inc.Status,
    inc.Title, inc.Description, inc.ReportedBy, inc.SLADueAt, inc.SLABreached, inc.CreatedAt, inc.UpdatedAt)
```

**After:**
```go
`, inc.ID, inc.TenantID, inc.SchoolID, inc.DeviceID,
    inc.SchoolName, inc.CountyID, inc.CountyName, inc.SubCountyID, inc.SubCountyName,
    inc.ContactName, inc.ContactPhone, inc.ContactEmail,
    inc.DeviceSerial, inc.DeviceAssetTag, inc.DeviceModelID, inc.DeviceMake, inc.DeviceModel, inc.DeviceCategory,
    inc.Category, inc.Severity, inc.Status,
    inc.Title, inc.Description, inc.ReportedBy, inc.SLADueAt, inc.SLABreached, inc.CreatedAt, inc.UpdatedAt)
```

### File Upload Security

**File:** `/home/pato/opt/ESSP/services/ims-api/internal/handlers/attachments.go`

**Observations:**
- Uses presigned URLs for direct upload/download to MinIO
- Object keys are generated server-side (prevents path traversal)
- Content-Type validation exists but could be stricter

**Recommendations:**
1. Add file type whitelist validation
2. Consider adding virus scanning for uploaded files
3. Implement file size validation at attachment creation time

---

## 6. Security Enhancements Implemented

### 6.1 Security Headers Middleware

**File:** `/home/pato/opt/ESSP/services/ims-api/internal/middleware/security.go`

**Implemented Headers:**

1. **X-Content-Type-Options: nosniff**
   - Prevents MIME type sniffing attacks
   - Forces browsers to respect declared content type

2. **X-Frame-Options: DENY**
   - Prevents clickjacking attacks
   - Blocks embedding in iframes

3. **X-XSS-Protection: 1; mode=block**
   - Enables XSS filter in older browsers
   - Blocks rendering on detected XSS

4. **Strict-Transport-Security: max-age=31536000; includeSubDomains**
   - Enforces HTTPS for 1 year
   - Includes all subdomains

5. **Content-Security-Policy: default-src 'none'**
   - Strict CSP appropriate for API
   - Prevents loading any external resources

6. **Referrer-Policy: strict-origin-when-cross-origin**
   - Controls referrer information leakage
   - Sends origin only on cross-origin requests

### 6.2 Request Body Size Limit Middleware

**File:** `/home/pato/opt/ESSP/services/ims-api/internal/middleware/security.go`

**Implementation:**
```go
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler
```

**Features:**
- Checks Content-Length header first for early rejection
- Uses `http.MaxBytesReader` to enforce limit during read
- Prevents memory exhaustion attacks
- Default limit: 10MB

**Applied in:** `/home/pato/opt/ESSP/services/ims-api/internal/api/server.go`

---

## 7. Server Configuration Review

### HTTP Server Settings

**File:** `/home/pato/opt/ESSP/services/ims-api/cmd/api/main.go`

**Good Practices:**
```go
httpServer := &http.Server{
    Addr:              cfg.HTTPAddr,
    Handler:           srv.Router(),
    ReadHeaderTimeout: 5 * time.Second,  // Prevents Slowloris attacks
}
```

**Recommendations:**
Consider adding:
```go
ReadTimeout:       10 * time.Second,
WriteTimeout:      10 * time.Second,
IdleTimeout:       120 * time.Second,
MaxHeaderBytes:    1 << 20, // 1MB
```

---

## 8. Middleware Stack Order

**Current Order (Correct):**
1. SecurityHeaders() - Applied first
2. MaxBodySize() - Size limit enforcement
3. RequestID() - Request tracking
4. Recoverer() - Panic recovery
5. Logger() - Request logging
6. CORS() - Cross-origin handling
7. AuthJWT() - Authentication (if enabled)
8. Tenancy() - Tenant/school context

**Security Note:** Security middleware is correctly applied before any request processing.

---

## 9. Remaining Concerns & Future Improvements

### Low Priority Items

1. **Rate Limiting**
   - No rate limiting middleware detected
   - Recommendation: Implement per-IP or per-tenant rate limiting
   - Suggested: Redis-based rate limiter

2. **API Versioning**
   - Routes are under `/v1` prefix (good)
   - Ensure backward compatibility when introducing `/v2`

3. **Audit Logging**
   - Current logging is operational
   - Recommendation: Add security audit log for sensitive operations
   - Track: login attempts, permission changes, data access

4. **Dependency Scanning**
   - Recommendation: Regular scanning of Go dependencies
   - Use: `go list -m all | nancy sleuth` or Dependabot

5. **Secret Rotation**
   - JWKS keys are refreshed (good)
   - Recommendation: Document database credential rotation process

6. **Database Connection Security**
   - Review: Ensure `sslmode=require` in production
   - Default shows: `sslmode=disable` (acceptable for dev only)

7. **Error Response Details**
   - Generic errors returned (good)
   - Ensure stack traces never leak to clients in production

---

## 10. Testing Recommendations

### Security Test Cases to Add

1. **SQL Injection Tests**
   - Test special characters in search queries
   - Test malicious input in all text fields

2. **Authentication Tests**
   - Test expired tokens
   - Test tokens with wrong audience
   - Test tokens with wrong issuer
   - Test "none" algorithm tokens (should fail)

3. **Authorization Tests**
   - Test cross-tenant data access
   - Test cross-school data access
   - Verify 404 vs 403 responses (prevent enumeration)

4. **Input Validation Tests**
   - Test oversized request bodies
   - Test extremely long strings
   - Test invalid JSON formats
   - Test negative numbers where not allowed

5. **CORS Tests**
   - Test cross-origin requests
   - Verify preflight handling

---

## 11. Compliance Considerations

### Data Protection

- **Tenant Isolation:** Fully implemented
- **Data Encryption in Transit:** HSTS enforced (when HTTPS enabled)
- **Data Encryption at Rest:** Database/storage layer responsibility
- **PII Handling:** Contact info stored (ensure compliance with data protection laws)

### Access Control

- **Authentication:** JWT-based (strong)
- **Authorization:** Tenant + School scoped (strong)
- **Audit Trail:** Basic operational logging (enhance for compliance)

---

## 12. Summary of Changes Made

### Code Fixes

1. **Fixed Critical Bug** in `internal/store/incidents_repo.go`
   - Corrected parameter count in INSERT statement
   - Added missing column values (school info, contact info, device info)

### New Files Created

1. **`internal/middleware/security.go`**
   - SecurityHeaders() middleware
   - MaxBodySize() middleware

### Modified Files

1. **`internal/api/server.go`**
   - Added SecurityHeaders() middleware
   - Added MaxBodySize(10MB) middleware
   - Applied before other middleware for proper security layering

### Documentation

1. **`docs/SECURITY_AUDIT.md`** (this file)
   - Complete security audit report
   - Findings and recommendations
   - Change log

---

## 13. Recommendations Summary

### Immediate (Already Implemented)
- [x] Fix incidents_repo.go parameter mismatch
- [x] Add security headers middleware
- [x] Add request body size limit
- [x] Document security audit findings

### Short Term (Next Sprint)
- [ ] Add HTTP server timeouts (ReadTimeout, WriteTimeout)
- [ ] Implement rate limiting middleware
- [ ] Add file type whitelist for uploads
- [ ] Add field length validation for string inputs

### Medium Term (Next Quarter)
- [ ] Implement security audit logging for sensitive operations
- [ ] Add comprehensive security test suite
- [ ] Set up automated dependency scanning
- [ ] Document secret rotation procedures

### Long Term (Next 6 Months)
- [ ] Consider adding virus scanning for file uploads
- [ ] Implement anomaly detection for unusual access patterns
- [ ] Add API request signing for extra critical endpoints
- [ ] Conduct penetration testing

---

## 14. Conclusion

The IMS API service demonstrates good security practices overall:

**Strengths:**
- No SQL injection vulnerabilities
- Strong authentication with JWT
- Proper tenant isolation throughout
- Good input validation practices
- Secrets properly externalized
- Security headers now implemented
- Request size limits now enforced

**Critical Issue Resolved:**
- incidents_repo.go parameter mismatch fixed

**Security Posture:** The service is production-ready from a security perspective with the implemented fixes. The recommendations provided are enhancements that should be prioritized based on risk assessment and business requirements.

**Next Review Date:** Recommended quarterly security audits or when major features are added.

---

**Report Prepared By:** Security Audit Team
**Date:** 2025-12-12
**Audit Reference:** SE-001
