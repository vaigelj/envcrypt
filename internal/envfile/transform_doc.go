// Package envfile provides utilities for parsing, writing, and manipulating
// .env files.
//
// The transform sub-feature allows bulk value transformations to be applied
// to a set of env entries. Built-in transform functions include:
//
//   - UppercaseValues: converts values to uppercase
//   - TrimValues: strips leading/trailing whitespace from values
//   - PrefixValues: prepends a fixed string to every value
//
// Custom TransformFunc implementations can be passed to Transform for
// project-specific needs.
package envfile
