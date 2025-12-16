package store

import (
	"encoding/base64"
	"strings"
	"time"
)

// Cursor format: base64("RFC3339Nano|id")
func EncodeCursor(t time.Time, id string) string {
	raw := t.UTC().Format(time.RFC3339Nano) + "|" + id
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func DecodeCursor(cur string) (time.Time, string, bool) {
	if cur == "" {
		return time.Time{}, "", false
	}
	b, err := base64.RawURLEncoding.DecodeString(cur)
	if err != nil {
		return time.Time{}, "", false
	}
	parts := strings.SplitN(string(b), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, "", false
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, "", false
	}
	return t, parts[1], true
}
