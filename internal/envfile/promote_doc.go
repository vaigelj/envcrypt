// Package envfile provides the Promote function for copying environment
// variables from one environment file (or set of entries) to another.
//
// Promote is useful when advancing configuration through deployment stages,
// for example from a development .env to a staging .env, while respecting
// existing values in the destination.
//
// Basic usage:
//
//	out, result, err := envfile.Promote(src, dst)
//
// Use WithPromoteOverwrite to allow destination values to be replaced, and
// WithPromoteExclude to skip specific keys during promotion.
//
// PromoteFile operates directly on file paths for convenience.
package envfile
