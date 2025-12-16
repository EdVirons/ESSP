package admin

import (
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// generateToken creates a new JWT token for the given username (legacy, admin only)
func (a *AdminAuth) generateToken(username string) (string, error) {
	return a.generateTokenForUser(&DemoUser{
		Username:    username,
		Roles:       []string{"ssp_admin"},
		Email:       "admin@essp.local",
		DisplayName: "Admin",
	})
}

// generateTokenForUser creates a new JWT token for a demo user with specific roles
func (a *AdminAuth) generateTokenForUser(user *DemoUser) (string, error) {
	now := time.Now()
	claims := AdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "essp-admin",
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(a.cfg.TokenExpiry)),
		},
		Username:    user.Username,
		Roles:       user.Roles,
		TenantID:    "default",
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(a.privateKey)
}

// validateToken validates a JWT token and returns its claims
func (a *AdminAuth) validateToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// setAuthCookie sets the authentication cookie
func (a *AdminAuth) setAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "essp_admin_token",
		Value:    token,
		Path:     "/",
		Domain:   a.cfg.CookieDomain,
		MaxAge:   int(a.cfg.TokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   a.cfg.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

// GetPublicKey returns the public key for token verification
// This can be used to expose a JWKS endpoint if needed
func (a *AdminAuth) GetPublicKey() *rsa.PublicKey {
	return a.publicKey
}
