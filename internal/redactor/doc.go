// Package redactor provides value-masking for sensitive environment variable
// keys before comparison results are rendered or reported.
//
// Keys are considered sensitive when their name contains substrings such as
// SECRET, PASSWORD, TOKEN, API_KEY, AUTH, PRIVATE, or CREDENTIAL (case-
// insensitive). Matching values in [comparator.Diff] entries are replaced
// with "***". Missing-key lists are never altered because they contain no
// values.
package redactor
