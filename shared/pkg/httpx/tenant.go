package httpx

import (
	"net/http"
	"strings"
)

func TenantID(r *http.Request) string {
	return strings.TrimSpace(r.Header.Get("X-Tenant-Id"))
}
