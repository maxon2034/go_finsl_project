package server

import (
	"testing"
)

func TestServerStart(t *testing.T) {

	err := ServerStart()
	if err != nil {
		t.Errorf("Server starting mistake")
	}
}
