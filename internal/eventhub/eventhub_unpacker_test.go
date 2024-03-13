package eventhub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshallEvent_NormalJSONInput(t *testing.T) {
	jsonInput := []byte(`{"records": [{"properties": {"log": "normal log entry"}}]}`)

	var event Event

	err := UnmarshallEvent(jsonInput, &event)

	assert.Nil(t, err)
}

func TestUnmarshallEvent_AbnormalJSONInput(t *testing.T) {
	// JSON with missing subproperty
	jsonInput := []byte(`{"records": [{"properties": {}}]}`)

	var event Event

	err := UnmarshallEvent(jsonInput, &event)

	assert.Nil(t, err)
}

func TestUnmarshallEvent_ExtraJSONInput(t *testing.T) {
	// JSON with extra property
	jsonInput := []byte(`{"records": [{"properties": {"log": "normal log entry"}, "extra": "additional property"}]}`)

	var event Event

	err := UnmarshallEvent(jsonInput, &event)

	assert.Nil(t, err)
}
