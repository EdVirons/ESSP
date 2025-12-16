package mocks_test

import (
	"testing"

	"github.com/edvirons/ssp/ims/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example test demonstrating MockNATSPublisher usage
func TestMockNATSPublisher_PublishJSON(t *testing.T) {
	publisher := mocks.NewMockNATSPublisher()

	// Publish a message
	type TestEvent struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	event := TestEvent{ID: "123", Name: "test"}
	err := publisher.PublishJSON("test.subject", event)
	require.NoError(t, err)

	// Verify message was published
	assert.True(t, publisher.WasPublished("test.subject"))
	assert.Equal(t, 1, publisher.GetMessageCount("test.subject"))

	// Get the message
	msg, ok := publisher.GetLastMessage("test.subject")
	require.True(t, ok)

	// Verify message content
	eventMsg, ok := msg.(TestEvent)
	require.True(t, ok)
	assert.Equal(t, "123", eventMsg.ID)
	assert.Equal(t, "test", eventMsg.Name)
}

func TestMockNATSPublisher_MultipleMessages(t *testing.T) {
	publisher := mocks.NewMockNATSPublisher()

	// Publish multiple messages
	for i := 1; i <= 5; i++ {
		err := publisher.PublishJSON("events", map[string]int{"count": i})
		require.NoError(t, err)
	}

	// Verify count
	assert.Equal(t, 5, publisher.GetMessageCount("events"))

	// Verify all messages
	messages := publisher.GetMessages("events")
	assert.Len(t, messages, 5)
}

func TestMockNATSPublisher_MultipleSubjects(t *testing.T) {
	publisher := mocks.NewMockNATSPublisher()

	// Publish to different subjects
	publisher.PublishJSON("subject1", "message1")
	publisher.PublishJSON("subject2", "message2")
	publisher.PublishJSON("subject1", "message3")

	// Verify subjects
	subjects := publisher.GetAllSubjects()
	assert.Len(t, subjects, 2)

	// Verify message counts
	assert.Equal(t, 2, publisher.GetMessageCount("subject1"))
	assert.Equal(t, 1, publisher.GetMessageCount("subject2"))
}

func TestMockNATSPublisher_AssertPublished(t *testing.T) {
	publisher := mocks.NewMockNATSPublisher()

	type Event struct {
		Type string `json:"type"`
		Data string `json:"data"`
	}

	// Publish some events
	publisher.PublishJSON("events", Event{Type: "create", Data: "item1"})
	publisher.PublishJSON("events", Event{Type: "update", Data: "item2"})
	publisher.PublishJSON("events", Event{Type: "delete", Data: "item3"})

	// Assert specific event was published
	found := publisher.AssertPublished("events", func(msg interface{}) bool {
		event, ok := msg.(Event)
		return ok && event.Type == "update" && event.Data == "item2"
	})
	assert.True(t, found)

	// Assert non-existent event
	found = publisher.AssertPublished("events", func(msg interface{}) bool {
		event, ok := msg.(Event)
		return ok && event.Type == "invalid"
	})
	assert.False(t, found)
}

func TestMockNATSPublisher_Reset(t *testing.T) {
	publisher := mocks.NewMockNATSPublisher()

	// Publish messages
	publisher.PublishJSON("test", "message1")
	publisher.PublishJSON("test", "message2")

	// Reset
	publisher.Reset()

	// Verify cleared
	assert.Equal(t, 0, publisher.GetMessageCount("test"))
	assert.False(t, publisher.WasPublished("test"))
}
