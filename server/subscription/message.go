package subscription

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	ResourceID string
	Type       string
	Content    interface{}
}

func (m Message) DecodeContent(output interface{}) error {
	return json.Unmarshal(m.Content.([]byte), output)
}

func (m Message) EncodeContent() (Message, error) {
	encoded, err := json.Marshal(m.Content)
	if err != nil {
		return Message{}, fmt.Errorf("failed to encode content: %w", err)
	}

	return Message{
		ResourceID: m.ResourceID,
		Type:       m.Type,
		Content:    encoded,
	}, nil
}

func (m Message) Encode() ([]byte, error) {
	// make sure the content is already encoded
	if _, ok := m.Content.([]byte); !ok {
		var err error
		m, err = m.EncodeContent()
		if err != nil {
			return nil, fmt.Errorf("failed to encode Message content: %w", err)
		}
	}

	return json.Marshal(m)
}

func DecodeMessage(data []byte) (Message, error) {
	var m Message
	if err := json.Unmarshal(data, &m); err != nil {
		return Message{}, fmt.Errorf("failed to decode message: %w", err)
	}

	return m, nil
}

