package eventhub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshallEvent_NormalJSONInput_ValidBehavior(t *testing.T) {
	jsonInput := []byte(`{"records": [{"properties": {"log": "normal log entry"}}]}`)

	err, _ := UnmarshallEvent(jsonInput)

	assert.Nil(t, err)
}

func TestUnmarshallEvent_MissingSubpropertyJSONInput_ValidBehavior(t *testing.T) {
	jsonInput := []byte(`{"records": [{"properties": {}}]}`)

	err, _ := UnmarshallEvent(jsonInput)

	assert.Nil(t, err)
}

func TestUnmarshallEvent_MissingEndQuoteJSONInput_ValidBehavior(t *testing.T) {
	jsonInput := []byte(`{"records": [{"properties": {"log": "missing closing quotation mark}}]}`)

	err, _ := UnmarshallEvent(jsonInput)

	assert.NotNil(t, err)
}

func TestUnmarshallEvent_ExtraPropertyJSONInput_ValidBehavior(t *testing.T) {
	jsonInput := []byte(
		`{"records": [{"properties": {"log": "normal log entry"}, "extra": "additional property"}]}`,
	)

	err, _ := UnmarshallEvent(jsonInput)

	assert.Nil(t, err)
}
