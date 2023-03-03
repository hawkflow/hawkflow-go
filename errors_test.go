package hawkflow

import (
	"testing"
)

func TestCreateError(t *testing.T) {
	err := createError("Test message.")
	expectedMessage := "Test message. Please see documentation at https://docs.hawkflow.ai/integration/index.html"

	if err.Error() != expectedMessage {
		t.Errorf("%v expected, got %v", expectedMessage, err.Error())
	}
}
