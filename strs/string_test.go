package strs

import "testing"

func TestStrOr(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		fallbacks []string
		want      string
	}{
		{
			name: "keeps original value",
			str:  "primary",
			fallbacks: []string{
				"fallback",
			},
			want: "primary",
		},
		{
			name: "uses first non-empty fallback",
			fallbacks: []string{
				"",
				"fallback",
				"later",
			},
			want: "fallback",
		},
		{
			name: "returns empty when no value exists",
			fallbacks: []string{
				"",
			},
			want: "",
		},
		{
			name: "handles no fallbacks",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrOr(tt.str, tt.fallbacks...); got != tt.want {
				t.Errorf("StrOr(%q, %q) = %q, want %q", tt.str, tt.fallbacks, got, tt.want)
			}
		})
	}
}

func TestStrIfElse(t *testing.T) {
	tests := []struct {
		ok      bool
		onTrue  string
		onFalse string
		want    string
	}{
		{ok: true, onTrue: "yes", onFalse: "no", want: "yes"},
		{ok: false, onTrue: "yes", onFalse: "no", want: "no"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := StrIfElse(tt.ok, tt.onTrue, tt.onFalse); got != tt.want {
				t.Errorf("StrIfElse(%t, %q, %q) = %q, want %q", tt.ok, tt.onTrue, tt.onFalse, got, tt.want)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		str    string
		prefix string
		want   bool
	}{
		{str: "gopher", prefix: "go", want: true},
		{str: "gopher", prefix: "ph", want: false},
		{str: "go", prefix: "gopher", want: false},
		{str: "gopher", prefix: "", want: true},
		{str: "", prefix: "", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.str+"/"+tt.prefix, func(t *testing.T) {
			if got := StartsWith(tt.str, tt.prefix); got != tt.want {
				t.Errorf("StartsWith(%q, %q) = %t, want %t", tt.str, tt.prefix, got, tt.want)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		str    string
		suffix string
		want   bool
	}{
		{str: "gopher", suffix: "her", want: true},
		{str: "gopher", suffix: "ph", want: false},
		{str: "go", suffix: "gopher", want: false},
		{str: "gopher", suffix: "", want: true},
		{str: "", suffix: "", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.str+"/"+tt.suffix, func(t *testing.T) {
			if got := EndsWith(tt.str, tt.suffix); got != tt.want {
				t.Errorf("EndsWith(%q, %q) = %t, want %t", tt.str, tt.suffix, got, tt.want)
			}
		})
	}
}

func TestPtr(t *testing.T) {
	input := "value"

	got := Ptr(input)
	if got == nil {
		t.Fatalf("Ptr(%q) = nil, want pointer to %q", input, input)
	}
	if *got != input {
		t.Errorf("Ptr(%q) points to %q, want %q", input, *got, input)
	}

	input = "changed"
	if *got != "value" {
		t.Errorf("Ptr(%q) changed after input reassignment: got %q, want %q", "value", *got, "value")
	}
}
