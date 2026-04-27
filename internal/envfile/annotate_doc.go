// Package envfile provides the annotate feature for attaching human-readable
// metadata to individual environment variable keys.
//
// # Overview
//
// Annotations are stored as a JSON file (.envcrypt_annotations.json) in the
// working directory alongside the .env file. Each annotation can carry:
//
//   - Description: a free-text explanation of the key's purpose
//   - Owner: the team or individual responsible for the key
//   - Deprecated: a flag indicating the key should no longer be used
//   - Tags: arbitrary labels for grouping or filtering
//
// # Usage
//
//	err := envfile.SetAnnotation(dir, "DATABASE_URL", envfile.Annotation{
//	    Description: "Primary Postgres connection string",
//	    Owner:       "platform-team",
//	})
//
//	ann, ok, err := envfile.GetAnnotation(dir, "DATABASE_URL")
package envfile
