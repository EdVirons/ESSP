package ws

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, validate origin against allowed origins
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub    *Hub
	logger *zap.Logger
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub, logger *zap.Logger) *Handler {
	return &Handler{
		hub:    hub,
		logger: logger,
	}
}

// ServeWS handles websocket requests from clients
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Extract tenant and user from headers or query params
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		tenantID = r.URL.Query().Get("tenant")
	}
	if tenantID == "" {
		tenantID = "default"
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = r.URL.Query().Get("user")
	}
	if userID == "" {
		// Try to extract from Authorization header (JWT)
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			// In production, decode JWT to get user ID
			userID = "anonymous"
		} else {
			userID = "anonymous"
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("failed to upgrade websocket", zap.Error(err))
		return
	}

	client := NewClient(h.hub, conn, userID, tenantID, h.logger)
	h.hub.register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}

// Hub returns the WebSocket hub
func (h *Handler) Hub() *Hub {
	return h.hub
}
