// Package envfile provides utilities for parsing, writing, and managing
// .env files.
//
// The tags sub-feature allows users to group environment variable keys
// under named labels (tags). Tags are persisted as a JSON sidecar file
// (.envcrypt_tags.json) alongside the env file directory, enabling
// workflows such as filtering exports or redaction by logical group.
package envfile
