package envfile

import "strings"

// SensitiveKeyPrefixes contains common prefixes for sensitive env var names.
var SensitiveKeyPrefixes = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"DSN",
	"DATABASE_URL",
}

// IsSensitive returns true if the key looks like it holds sensitive data.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, prefix := range SensitiveKeyPrefixes {
		if strings.Contains(upper, prefix) {
			return true
		}
	}
	return false
}

// Redact returns a copy of the entries map where sensitive values are masked.
// The mask string is used as the replacement value; if empty, "***" is used.
func Redact(entries map[string]string, mask string) map[string]string {
	if mask == "" {
		mask = "***"
	}
	out := make(map[string]string, len(entries))
	for k, v := range entries {
		if IsSensitive(k) {
			out[k] = mask
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactString replaces the value portion of a KEY=VALUE line if the key is
// sensitive. Lines that are comments or blank are returned unchanged.
func RedactString(line, mask string) string {
	if mask == "" {
		mask = "***"
	}
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return line
	}
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return line
	}
	key := strings.TrimSpace(line[:idx])
	if IsSensitive(key) {
		return key + "=" + mask
	}
	return line
}
