package llm

import "testing"

func TestValidateURL(t *testing.T) {
	c := NewClient()
	if err := c.ValidateURL("not a url", []string{}); err == nil {
		t.Fatal("expected error for invalid url")
	}
	if err := c.ValidateURL("ftp://example.com", []string{"example.com"}); err == nil {
		t.Fatal("expected error for non-http scheme")
	}
	if err := c.ValidateURL("https://api.openai.com", []string{"api.openai.com"}); err != nil {
		t.Fatalf("expected ok, got %v", err)
	}
	if err := c.ValidateURL("https://blocked.com", []string{"api.openai.com"}); err == nil {
		t.Fatal("expected blocked host")
	}
}
