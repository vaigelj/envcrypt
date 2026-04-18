package envfile

// MergeStrategy controls how conflicts are resolved when merging two env maps.
type MergeStrategy int

const (
	// PreferBase keeps the value from the base map on conflict.
	PreferBase MergeStrategy = iota
	// PreferOverride replaces base values with override values on conflict.
	PreferOverride
)

// Merge combines base and override maps according to the given strategy.
// Keys present only in override are always added to the result.
func Merge(base, override map[string]string, strategy MergeStrategy) map[string]string {
	result := make(map[string]string, len(base))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		if _, exists := result[k]; !exists {
			result[k] = v
			continue
		}
		if strategy == PreferOverride {
			result[k] = v
		}
	}
	return result
}

// Diff returns keys whose values differ between a and b, plus keys present in
// only one of the maps. The returned map contains the value from b (the "new"
// side) for each differing or added key. Removed keys map to an empty string.
func Diff(a, b map[string]string) map[string]string {
	diff := make(map[string]string)
	for k, va := range a {
		if vb, ok := b[k]; !ok {
			diff[k] = ""
		} else if va != vb {
			diff[k] = vb
		}
	}
	for k, vb := range b {
		if _, ok := a[k]; !ok {
			diff[k] = vb
		}
	}
	return diff
}
