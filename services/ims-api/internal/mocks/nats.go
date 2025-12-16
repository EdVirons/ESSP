package mocks

import (
	"encoding/json"
	"fmt"
	"sync"
)

// MockNATSPublisher is a mock implementation of NATS publisher for testing.
type MockNATSPublisher struct {
	mu       sync.RWMutex
	messages map[string][]interface{} // subject -> list of published messages
	calls    map[string]int
	errors   map[string]error
}

// NewMockNATSPublisher creates a new mock NATS publisher.
func NewMockNATSPublisher() *MockNATSPublisher {
	return &MockNATSPublisher{
		messages: make(map[string][]interface{}),
		calls:    make(map[string]int),
		errors:   make(map[string]error),
	}
}

// PublishJSON publishes a JSON message to a subject.
func (m *MockNATSPublisher) PublishJSON(subject string, v interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("PublishJSON")

	if err := m.errors["PublishJSON"]; err != nil {
		return err
	}

	// Validate that the message can be marshaled to JSON
	_, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Store the message
	m.messages[subject] = append(m.messages[subject], v)
	return nil
}

// Publish publishes a raw byte message to a subject.
func (m *MockNATSPublisher) Publish(subject string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Publish")

	if err := m.errors["Publish"]; err != nil {
		return err
	}

	m.messages[subject] = append(m.messages[subject], data)
	return nil
}

// GetMessages returns all messages published to a subject.
func (m *MockNATSPublisher) GetMessages(subject string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, exists := m.messages[subject]
	if !exists {
		return []interface{}{}
	}

	// Return a copy to avoid race conditions
	result := make([]interface{}, len(messages))
	copy(result, messages)
	return result
}

// GetMessageCount returns the number of messages published to a subject.
func (m *MockNATSPublisher) GetMessageCount(subject string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.messages[subject])
}

// GetLastMessage returns the last message published to a subject.
func (m *MockNATSPublisher) GetLastMessage(subject string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, exists := m.messages[subject]
	if !exists || len(messages) == 0 {
		return nil, false
	}

	return messages[len(messages)-1], true
}

// WasPublished checks if any message was published to a subject.
func (m *MockNATSPublisher) WasPublished(subject string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.messages[subject]) > 0
}

// GetAllSubjects returns all subjects that have received messages.
func (m *MockNATSPublisher) GetAllSubjects() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	subjects := make([]string, 0, len(m.messages))
	for subject := range m.messages {
		subjects = append(subjects, subject)
	}
	return subjects
}

// Reset clears all stored messages and call counts.
func (m *MockNATSPublisher) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = make(map[string][]interface{})
	m.calls = make(map[string]int)
	m.errors = make(map[string]error)
}

// SetError configures the mock to return an error for a specific method.
func (m *MockNATSPublisher) SetError(method string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors[method] = err
}

// ClearErrors removes all configured errors.
func (m *MockNATSPublisher) ClearErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors = make(map[string]error)
}

// GetCallCount returns the number of times a method was called.
func (m *MockNATSPublisher) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[method]
}

// ClearSubject removes all messages for a specific subject.
func (m *MockNATSPublisher) ClearSubject(subject string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.messages, subject)
}

// AssertPublished is a helper for testing that a message was published.
// Returns true if a message matching the predicate was found.
func (m *MockNATSPublisher) AssertPublished(subject string, predicate func(interface{}) bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages, exists := m.messages[subject]
	if !exists {
		return false
	}

	for _, msg := range messages {
		if predicate(msg) {
			return true
		}
	}
	return false
}

// trackCall increments the call counter for a method (must be called with lock held).
func (m *MockNATSPublisher) trackCall(method string) {
	m.calls[method]++
}
