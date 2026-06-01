package services

import "testing"

func TestMaskKey(t *testing.T) {
	if got := mask(""); got != "****" {
		t.Errorf("expected ****, got %s", got)
	}
	if got := mask("short"); got != "****" {
		t.Errorf("expected ****, got %s", got)
	}
	if got := mask("sk-abcdefghijklmnop"); got != "sk-a...mnop" {
		t.Errorf("got %s", got)
	}
}
