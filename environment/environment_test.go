package environment

import "testing"

func TestGet(t *testing.T) {
	_, err := Initialize()
	if err != nil {
		t.Errorf("GetEnvironment() failed: %v", err)
	}
}
