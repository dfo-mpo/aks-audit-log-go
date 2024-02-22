package eventhub

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	// Call the function with a specific size
	size := 10
	result, err := generate(size)

	// Assert the result and error
	if err != nil {
		t.Errorf("generate returned an error: %v", err)
	}

	// Assert the length of the generated string
	if len(result) != size {
		t.Errorf("generate did not return a string of the expected size. Expected size: %d, Actual size: %d", size, len(result))
	}

	// Assert that the generated string contains only characters from the alphabet
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, char := range result {
		if strings.ContainsRune(alphabet, char) {
			t.Errorf("generate returned a string containing an invalid character: %c", char)
		}
	}
}
