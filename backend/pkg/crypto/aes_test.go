package crypto

import "testing"

func TestRoundTrip(t *testing.T) {
	c, err := New("01234567890123456789012345678901")
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	ct, err := c.Encrypt("hello world")
	if err != nil {
		t.Fatalf("enc: %v", err)
	}
	pt, err := c.Decrypt(ct)
	if err != nil {
		t.Fatalf("dec: %v", err)
	}
	if pt != "hello world" {
		t.Fatalf("expected hello world, got %s", pt)
	}
}

func TestEmpty(t *testing.T) {
	c, _ := New("01234567890123456789012345678901")
	ct, _ := c.Encrypt("")
	if ct != "" {
		t.Fatal("expected empty ciphertext for empty input")
	}
	pt, _ := c.Decrypt("")
	if pt != "" {
		t.Fatal("expected empty plaintext for empty input")
	}
}

func TestInvalidKey(t *testing.T) {
	if _, err := New("short"); err == nil {
		t.Fatal("expected error for short key")
	}
}
