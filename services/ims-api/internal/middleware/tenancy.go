package middleware

import (
	"net/http"
	"strings"

	"github.com/edvirons/ssp/ims/internal/config"
)

func Tenancy(cfg config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenant := strings.TrimSpace(r.Header.Get(cfg.TenantHeader))
			school := strings.TrimSpace(r.Header.Get(cfg.SchoolHeader))

			if tenant == "" {
				tenant = cfg.DevTenantID
			}
			if school == "" {
				school = cfg.DevSchoolID
			}

			ctx := WithTenantID(r.Context(), tenant)
			ctx = WithSchoolID(ctx, school)

			// In dev mode (auth disabled), assign admin role for full access
			if !cfg.AuthEnabled {
				ctx = WithRoles(ctx, []string{"ssp_admin"})
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
