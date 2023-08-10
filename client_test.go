package coprime

import (
	"errors"
	"testing"
)

func TestClientErrorsOnNotFound(t *testing.T) {
	client := NewTestClient()
	_, err := client.Request("GET", Sandbox, "/fake", nil, nil)
	if err == nil {
		t.Error(errors.New("Should have thrown 404 error"))
	}
}
