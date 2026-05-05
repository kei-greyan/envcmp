package comparator

// Result holds the comparison outcome between two env files.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Mismatched contains keys present in both files but with different values.
	Mismatched []MismatchedKey
}

// MismatchedKey describes a key whose value differs between two env files.
type MismatchedKey struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare performs a diff between two parsed env maps.
// left and right are maps of key->value as returned by parser.ParseFile.
func Compare(left, right map[string]string) Result {
	result := Result{}

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			result.MissingInRight = append(result.MissingInRight, k)
			continue
		}
		if lv != rv {
			result.Mismatched = append(result.Mismatched, MismatchedKey{
				Key:        k,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for k := range right {
		if _, ok := left[k]; !ok {
			result.MissingInLeft = append(result.MissingInLeft, k)
		}
	}

	return result
}

// HasDiff returns true when the Result contains any differences.
func (r Result) HasDiff() bool {
	return len(r.MissingInRight) > 0 || len(r.MissingInLeft) > 0 || len(r.Mismatched) > 0
}
