package analysis

import "testing"

func TestDeriveRisk(t *testing.T) {
	cases := []struct {
		name string
		risks []any
		want string
	}{
		{"empty", nil, "low"},
		{"single", []any{map[string]any{"level": "high"}}, "high"},
		{"max", []any{map[string]any{"level": "low"}, map[string]any{"level": "critical"}}, "critical"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := &Schema{Risks: c.risks}
			if got := deriveRisk(s); got != c.want {
				t.Fatalf("got %s want %s", got, c.want)
			}
		})
	}
}

func TestBuildUserPrompt(t *testing.T) {
	p := BuildUserPrompt("foo", "procedure", "BEGIN NULL; END;")
	if p == "" {
		t.Fatal("expected prompt")
	}
}
