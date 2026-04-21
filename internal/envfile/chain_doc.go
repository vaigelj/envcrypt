// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// The Chain type enables layered configuration by merging an ordered list of
// .env files. Entries from later files override those from earlier files,
// following the same semantics as Merge with preferOverride=true.
//
// Files that do not exist on disk are silently skipped, making it safe to
// include optional override files (e.g. ".env.local") in the chain without
// requiring them to be present.
//
// Typical usage:
//
//	chain := envfile.NewChain(".env", ".env.local", ".env.production")
//	env, err := chain.ResolveMap()
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The resolved map can then be used to populate process environment variables,
// pass configuration to subsystems, or feed into an encryption/decryption
// pipeline provided by the envcrypt package.
package envfile
