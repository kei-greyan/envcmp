// Package ignore provides support for loading and applying ignore rules
// to comparison results. An ignore file is a plain text file where each
// line specifies a key name that should be excluded from the diff output.
//
// Lines beginning with '#' are treated as comments and are skipped,
// as are blank lines.
//
// Example ignore file:
//
//	# credentials — always differ
//	AWS_SECRET_ACCESS_KEY
//	AWS_ACCESS_KEY_ID
//
//	# local-only overrides
//	PORT
package ignore
