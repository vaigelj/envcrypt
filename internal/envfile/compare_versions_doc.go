// Package envfile provides compare_versions utilities for diffing two named
// env file versions. Use CompareVersions to produce a VersionDiff, and
// FormatVersionDiff to render it as a human-readable string.
//
// Example:
//
//	from := envfile.Version{Name: "staging", Vars: map[string]string{"PORT": "8080"}}
//	to   := envfile.Version{Name: "prod",    Vars: map[string]string{"PORT": "443"}}
//	diff := envfile.CompareVersions(from, to)
//	fmt.Print(envfile.FormatVersionDiff(diff))
package envfile
