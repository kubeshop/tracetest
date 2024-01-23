package subscription_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestDecode(t *testing.T) {
	type msg struct {
		Value        int
		AnotherValue string
	}

	realMessage := msg{
		Value:        2,
		AnotherValue: "cat",
	}

	message := subscription.Message{
		ResourceID: "xxx",
		Content:    realMessage,
	}

	msgBytes, err := json.Marshal(message)
	require.NoError(t, err)

	var target subscription.Message
	err = json.Unmarshal(msgBytes, &target)
	require.NoError(t, err)

	var targetMsg msg
	err = target.DecodeContent(&targetMsg)
	require.NoError(t, err)

	assert.Equal(t, 2, targetMsg.Value)
	assert.Equal(t, "cat", targetMsg.AnotherValue)
}
