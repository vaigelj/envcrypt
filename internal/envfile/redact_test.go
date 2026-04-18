package envfile

import "testing"

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key       string
		wantSensitive bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"SECRET_TOKEN", true},
		{"AUTH_HEADER", true},
		{"DATABASE_URL", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"LOG_LEVEL", false},
	}
	for _, tc := range cases {
		got := IsSensitive(tc.key)
		if got != tc.wantSensitive {
			t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.wantSensitive)
		}
	}
}

func TestRedact(t *testing.T) {
	input := map[string]string{
		"APP_NAME":    "myapp",
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"PORT":        "8080",
	}
	out := Redact(input, "")
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT unchanged, got %q", out["PORT"])
	}
	if out["DB_PASSWORD"] != "***" {
		t.Errorf("expected DB_PASSWORD redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "***" {
		t.Errorf("expected API_KEY redacted, got %q", out["API_KEY"])
	}
}

func TestRedactCustomMask(t *testing.T) {
	input := map[string]string{"SECRET_KEY": "value"}
	out := Redact(input, "<hidden>")
	if out["SECRET_KEY"] != "<hidden>" {
		t.Errorf("expected custom mask, got %q", out["SECRET_KEY"])
	}
}

func TestRedactString(t *testing.T) {
	cases := []struct {
		line string
		want string
	}{
		{"APP_NAME=myapp", "APP_NAME=myapp"},
		{"DB_PASSWORD=secret", "DB_PASSWORD=***"},
		{"# comment line", "# comment line"},
		{"", ""},
		{"INVALID_LINE", "INVALID_LINE"},
		{"API_KEY=tok123", "API_KEY=***"},
	}
	for _, tc := range cases {
		got := RedactString(tc.line, "")
		if got != tc.want {
			t.Errorf("RedactString(%q) = %q, want %q", tc.line, got, tc.want)
		}
	}
}
