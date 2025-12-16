package ws

import (
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	MessageTypeNotification      MessageType = "notification"
	MessageTypeEntityUpdate      MessageType = "entity_update"
	MessageTypePing              MessageType = "ping"
	MessageTypePong              MessageType = "pong"
	MessageTypeChatMessage       MessageType = "chat_message"
	MessageTypeChatTyping        MessageType = "chat_typing"
	MessageTypeChatRead          MessageType = "chat_read"
	MessageTypeChatSessionUpdate MessageType = "chat_session_update"
	MessageTypePresenceUpdate    MessageType = "presence_update"
	MessageTypeNewChatWaiting    MessageType = "new_chat_waiting"
	MessageTypeAgentAssigned     MessageType = "agent_assigned"
	MessageTypeTypingIndicator   MessageType = "typing_indicator"
)

// Message is a WebSocket message
type Message struct {
	Type      MessageType `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp string      `json:"timestamp"`
}

// NotificationPayload is the payload for notification messages
type NotificationPayload struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Action    string                 `json:"action"`
	Actor     string                 `json:"actor"`
	Target    string                 `json:"target"`
	Summary   string                 `json:"summary"`
	Timestamp string                 `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	// Registered clients by tenant
	clients map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to all clients
	broadcast chan *BroadcastMessage

	// Logger
	logger *zap.Logger

	// Mutex for clients map
	mu sync.RWMutex
}

// BroadcastMessage is a message to broadcast to clients
type BroadcastMessage struct {
	TenantID string
	Message  *Message
}

// NewHub creates a new Hub
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
		logger:     logger,
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.tenantID] == nil {
				h.clients[client.tenantID] = make(map[*Client]bool)
			}
			h.clients[client.tenantID][client] = true
			h.mu.Unlock()
			h.logger.Debug("client registered",
				zap.String("tenantId", client.tenantID),
				zap.String("userId", client.userID),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.tenantID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.tenantID)
					}
				}
			}
			h.mu.Unlock()
			h.logger.Debug("client unregistered",
				zap.String("tenantId", client.tenantID),
				zap.String("userId", client.userID),
			)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[msg.TenantID]
			h.mu.RUnlock()

			data, err := json.Marshal(msg.Message)
			if err != nil {
				h.logger.Error("failed to marshal message", zap.Error(err))
				continue
			}

			for client := range clients {
				select {
				case client.send <- data:
				default:
					// Client buffer full, close connection
					h.mu.Lock()
					close(client.send)
					delete(h.clients[msg.TenantID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

// Broadcast sends a message to all clients of a tenant
func (h *Hub) Broadcast(tenantID string, msg *Message) {
	if msg.Timestamp == "" {
		msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	h.broadcast <- &BroadcastMessage{
		TenantID: tenantID,
		Message:  msg,
	}
}

// BroadcastNotification sends a notification to all clients of a tenant
func (h *Hub) BroadcastNotification(tenantID string, payload NotificationPayload) {
	h.Broadcast(tenantID, &Message{
		Type:    MessageTypeNotification,
		Payload: payload,
	})
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	count := 0
	for _, clients := range h.clients {
		count += len(clients)
	}
	return count
}

// TenantClientCount returns the number of connected clients for a tenant
func (h *Hub) TenantClientCount(tenantID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[tenantID])
}
