package services

import "testing"

func TestAuthServiceSentinel(t *testing.T) {
	if ErrInvalidCredentials == nil || ErrInvalidToken == nil {
		t.Fatal("expected sentinels")
	}
}
