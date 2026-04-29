package envfile

import (
	"fmt"
	"strconv"
	"strings"
)

// CastOption configures a Cast operation.
type CastOption func(*castConfig)

type castConfig struct {
	keys    map[string]string // key -> target type
	strict  bool
}

// WithCastKeys specifies which keys to cast and to what type.
// Supported types: "int", "float", "bool", "string".
func WithCastKeys(rules map[string]string) CastOption {
	return func(c *castConfig) {
		for k, v := range rules {
			c.keys[k] = strings.ToLower(v)
		}
	}
}

// WithCastStrict causes Cast to return an error if a value cannot be
// converted to the requested type instead of leaving it unchanged.
func WithCastStrict() CastOption {
	return func(c *castConfig) { c.strict = true }
}

// Cast normalises the string values of env entries according to type rules.
// For example, "true" / "1" / "yes" all become "true" when cast to bool.
func Cast(entries []Entry, opts ...CastOption) ([]Entry, error) {
	cfg := &castConfig{keys: make(map[string]string)}
	for _, o := range opts {
		o(cfg)
	}

	out := make([]Entry, len(entries))
	copy(out, entries)

	for i, e := range out {
		typ, ok := cfg.keys[e.Key]
		if !ok {
			continue
		}
		casted, err := castValue(e.Value, typ)
		if err != nil {
			if cfg.strict {
				return nil, fmt.Errorf("cast: key %q value %q: %w", e.Key, e.Value, err)
			}
			continue
		}
		out[i].Value = casted
	}
	return out, nil
}

func castValue(val, typ string) (string, error) {
	switch typ {
	case "int":
		v, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to int", val)
		}
		return strconv.FormatInt(int64(v), 10), nil
	case "float":
		_, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return "", fmt.Errorf("cannot cast %q to float", val)
		}
		return strings.TrimSpace(val), nil
	case "bool":
		norm := strings.ToLower(strings.TrimSpace(val))
		switch norm {
		case "1", "true", "yes", "on":
			return "true", nil
		case "0", "false", "no", "off", "":
			return "false", nil
		default:
			return "", fmt.Errorf("cannot cast %q to bool", val)
		}
	case "string":
		return val, nil
	default:
		return "", fmt.Errorf("unknown type %q", typ)
	}
}
