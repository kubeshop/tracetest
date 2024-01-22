package subscription

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Message struct {
	ResourceID string
	Type       string
	Content    interface{}
}

func (m Message) DecodeContent(output interface{}) error {
	if _, isString := m.Content.(string); isString {
		base64Decoded, err := base64.StdEncoding.DecodeString(m.Content.(string))
		if err != nil {
			return fmt.Errorf("failed to decode base64 string: %w", err)
		}
		m.Content = base64Decoded
	}

	if _, isBytes := m.Content.(string); !isBytes {
		bytes, err := json.Marshal(m.Content)
		if err != nil {
			return fmt.Errorf("could not marshal json: %w", err)
		}
		m.Content = bytes
	}

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
