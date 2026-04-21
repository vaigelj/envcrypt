// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// The Chain type enables layered configuration by merging an ordered list of
// .env files. Entries from later files override those from earlier files,
// following the same semantics as Merge with preferOverride=true.
//
// Typical usage:
//
//	chain := envfile.NewChain(".env", ".env.local", ".env.production")
//	env, err := chain.ResolveMap()
package envfile
