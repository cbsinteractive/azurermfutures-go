package azurermfutures

import "testing"

func TestNewClient(t *testing.T) {
	_, err := New(nil, SetHost(""), SetAuthHost(""))
	if err != nil {
		t.Fatal(err)
	}
}
