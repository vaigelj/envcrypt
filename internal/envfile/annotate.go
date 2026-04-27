package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Annotation holds metadata attached to a specific env key.
type Annotation struct {
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Annotations maps env key names to their Annotation.
type Annotations map[string]Annotation

func annotationsPath(dir string) string {
	return filepath.Join(dir, ".envcrypt_annotations.json")
}

// LoadAnnotations reads annotations from dir. Returns empty map if missing.
func LoadAnnotations(dir string) (Annotations, error) {
	path := annotationsPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Annotations{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load annotations: %w", err)
	}
	var a Annotations
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("parse annotations: %w", err)
	}
	return a, nil
}

// SaveAnnotations writes annotations to dir.
func SaveAnnotations(dir string, a Annotations) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal annotations: %w", err)
	}
	return os.WriteFile(annotationsPath(dir), data, 0o600)
}

// SetAnnotation adds or updates the annotation for key in dir.
func SetAnnotation(dir, key string, ann Annotation) error {
	a, err := LoadAnnotations(dir)
	if err != nil {
		return err
	}
	a[key] = ann
	return SaveAnnotations(dir, a)
}

// GetAnnotation returns the annotation for key, and whether it existed.
func GetAnnotation(dir, key string) (Annotation, bool, error) {
	a, err := LoadAnnotations(dir)
	if err != nil {
		return Annotation{}, false, err
	}
	ann, ok := a[key]
	return ann, ok, nil
}

// RemoveAnnotation deletes the annotation for key from dir.
func RemoveAnnotation(dir, key string) error {
	a, err := LoadAnnotations(dir)
	if err != nil {
		return err
	}
	delete(a, key)
	return SaveAnnotations(dir, a)
}
