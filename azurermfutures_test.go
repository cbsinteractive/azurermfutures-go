package azurermfuturesgo

import "testing"

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Error("NewClient failed to return a value")
	}
}
