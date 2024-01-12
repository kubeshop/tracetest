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

// func decodeTimeFunc(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
// 	// From https://github.com/mitchellh/mapstructure/issues/159#issuecomment-482201507
// 	if t != reflect.TypeOf(time.Time{}) {
// 		return data, nil
// 	}

// 	switch f.Kind() {
// 	case reflect.String:
// 		return time.Parse(time.RFC3339, data.(string))
// 	case reflect.Float64:
// 		return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
// 	case reflect.Int64:
// 		return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
// 	default:
// 		return data, nil
// 	}
// }

// func (m Message) DecodeContent(output interface{}) error {
// 	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
// 		Metadata: nil,
// 		DecodeHook: mapstructure.ComposeDecodeHookFunc(
// 			decodeTimeFunc,
// 		),
// 		Result: output,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	if err := decoder.Decode(m.Content); err != nil {
// 		return err
// 	}
// 	return err
// }
