// Package envfile provides utilities for parsing, writing, and managing
// .env files. The profile sub-feature allows teams to maintain multiple
// named environment profiles (e.g. .env.dev, .env.staging, .env.prod)
// within the same directory and load or save them by name.
//
// Usage:
//
//	profiles, err := envfile.ListProfiles("/path/to/project")
//	data, err   := envfile.LoadProfile("/path/to/project", "dev")
//	err          = envfile.SaveProfile("/path/to/project", "dev", data)
package envfile
