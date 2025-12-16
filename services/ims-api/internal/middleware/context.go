package middleware

import "context"

type ctxKey string

const (
	ctxTenantID        ctxKey = "tenantId"
	ctxSchoolID        ctxKey = "schoolId"
	ctxRoles           ctxKey = "roles"
	ctxAssignedSchools ctxKey = "assignedSchools"
	ctxClaims          ctxKey = "claims"
	ctxUserID          ctxKey = "userId"
	ctxUserName        ctxKey = "userName"
)

func WithTenantID(ctx context.Context, tenant string) context.Context {
	return context.WithValue(ctx, ctxTenantID, tenant)
}
func WithSchoolID(ctx context.Context, school string) context.Context {
	return context.WithValue(ctx, ctxSchoolID, school)
}
func TenantID(ctx context.Context) string {
	v, _ := ctx.Value(ctxTenantID).(string)
	return v
}
func SchoolID(ctx context.Context) string {
	v, _ := ctx.Value(ctxSchoolID).(string)
	return v
}

// WithRoles stores the user's roles in the context
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, ctxRoles, roles)
}

// Roles retrieves the user's roles from the context
func Roles(ctx context.Context) []string {
	v, _ := ctx.Value(ctxRoles).([]string)
	return v
}

// WithAssignedSchools stores the user's assigned schools in the context
func WithAssignedSchools(ctx context.Context, schools []string) context.Context {
	return context.WithValue(ctx, ctxAssignedSchools, schools)
}

// AssignedSchools retrieves the user's assigned schools from the context
func AssignedSchools(ctx context.Context) []string {
	v, _ := ctx.Value(ctxAssignedSchools).([]string)
	return v
}

// WithClaims stores the full JWT claims in the context
func WithClaims(ctx context.Context, claims map[string]any) context.Context {
	return context.WithValue(ctx, ctxClaims, claims)
}

// Claims retrieves the full JWT claims from the context
func Claims(ctx context.Context) map[string]any {
	v, _ := ctx.Value(ctxClaims).(map[string]any)
	return v
}

// WithUserID stores the user ID in the context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxUserID, userID)
}

// UserID retrieves the user ID from the context
func UserID(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserID).(string)
	return v
}

// WithUserName stores the user's display name in the context
func WithUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, ctxUserName, userName)
}

// UserName retrieves the user's display name from the context
func UserName(ctx context.Context) string {
	v, _ := ctx.Value(ctxUserName).(string)
	return v
}
