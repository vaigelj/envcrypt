// Package envfile provides utilities for managing .env files.
//
// # Convert
//
// The Convert sub-feature allows re-encoding a set of environment entries
// from one serialisation format to another without losing key/value data.
//
// Supported formats:
//
//	"dotenv" – classic KEY=VALUE lines
//	"json"   – JSON object { "KEY": "value" }
//	"shell"  – export KEY=VALUE lines suitable for shell sourcing
//
// Example:
//
//	out, err := envfile.ConvertFormat(entries, "dotenv", "json")
//
//	err := envfile.ConvertFile("prod.env", "prod.json", "dotenv", "json")
package envfile
