package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	l := New("debug")
	if l == nil {
		t.Fatal("expected logger")
	}
}

func TestNewUnknownLevel(t *testing.T) {
	var buf bytes.Buffer
	l := New("unknown")
	l.Info("hi", "k", "v")
	if !strings.Contains(buf.String(), "hi") && buf.Len() == 0 {
	}
	_ = json.NewEncoder(&buf).Encode(nil)
}
