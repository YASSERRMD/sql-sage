package llm

import "testing"

func TestTrimSlash(t *testing.T) {
	if trimSlash("https://x.com/v1/") != "https://x.com/v1" {
		t.Fatal("expected trimmed")
	}
	if trimSlash("https://x.com/v1") != "https://x.com/v1" {
		t.Fatal("no change")
	}
}

func TestTruncate(t *testing.T) {
	if truncate("abcdef", 3) != "abc..." {
		t.Fatalf("got %s", truncate("abcdef", 3))
	}
	if truncate("abc", 3) != "abc" {
		t.Fatalf("got %s", truncate("abc", 3))
	}
}
