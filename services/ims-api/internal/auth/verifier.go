package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	issuer   string
	jwksURL  string
	audience string

	mu        sync.RWMutex
	keys      map[string]*rsa.PublicKey
	fetchedAt time.Time
	client    *http.Client
}

func NewVerifier(issuer, jwksURL, audience string) *Verifier {
	return &Verifier{
		issuer:   strings.TrimRight(issuer, "/"),
		jwksURL:  jwksURL,
		audience: audience,
		keys:     map[string]*rsa.PublicKey{},
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (v *Verifier) Verify(ctx context.Context, tokenString string) (map[string]any, error) {
	if err := v.ensureKeys(ctx); err != nil {
		return nil, err
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
	tok, err := parser.Parse(tokenString, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("missing kid")
		}
		v.mu.RLock()
		k := v.keys[kid]
		v.mu.RUnlock()
		if k == nil {
			_ = v.refreshKeys(ctx)
			v.mu.RLock()
			k = v.keys[kid]
			v.mu.RUnlock()
		}
		if k == nil {
			return nil, errors.New("unknown kid")
		}
		return k, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid claims")
	}

	if iss, _ := claims["iss"].(string); iss != "" && v.issuer != "" && strings.TrimRight(iss, "/") != v.issuer {
		return nil, errors.New("issuer mismatch")
	}
	if v.audience != "" {
		aud, _ := claims.GetAudience()
		found := false
		for _, a := range aud {
			if a == v.audience {
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("audience mismatch")
		}
	}

	out := map[string]any{}
	for k, val := range claims {
		out[k] = val
	}
	return out, nil
}

func (v *Verifier) ensureKeys(ctx context.Context) error {
	v.mu.RLock()
	stale := time.Since(v.fetchedAt) > 10*time.Minute || len(v.keys) == 0
	v.mu.RUnlock()
	if stale {
		return v.refreshKeys(ctx)
	}
	return nil
}

type jwks struct {
	Keys []struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

func (v *Verifier) refreshKeys(ctx context.Context) error {
	if v.jwksURL == "" {
		return errors.New("jwks url not set")
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var j jwks
	if err := json.NewDecoder(resp.Body).Decode(&j); err != nil {
		return err
	}

	keys := map[string]*rsa.PublicKey{}
	for _, k := range j.Keys {
		if k.Kty != "RSA" || k.Kid == "" || k.N == "" || k.E == "" {
			continue
		}
		pub, err := rsaFromJWK(k.N, k.E)
		if err != nil {
			continue
		}
		keys[k.Kid] = pub
	}
	if len(keys) == 0 {
		return errors.New("no keys found")
	}

	v.mu.Lock()
	v.keys = keys
	v.fetchedAt = time.Now()
	v.mu.Unlock()
	return nil
}

func rsaFromJWK(nB64, eB64 string) (*rsa.PublicKey, error) {
	nb, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, err
	}
	eb, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, err
	}
	n := new(big.Int).SetBytes(nb)
	e := new(big.Int).SetBytes(eb).Int64()
	if e == 0 {
		return nil, errors.New("invalid exponent")
	}
	return &rsa.PublicKey{N: n, E: int(e)}, nil
}
