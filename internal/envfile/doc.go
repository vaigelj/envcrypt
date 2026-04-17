// Package envfile handles reading and writing .env files for the envcrypt tool.
//
// It supports the standard KEY=VALUE format, with optional double-quoted values,
// blank line skipping, and comment lines prefixed with '#'.
//
// Typical usage:
//
//	ef, err := envfile.Parse(".env")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, entry := range ef.Entries {
//		fmt.Println(entry.Key, "=", entry.Value)
//	}
package envfile
