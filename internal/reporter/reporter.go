package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/envcmp/internal/comparator"
)

// Report writes a human-readable diff summary to w.
// leftName and rightName are labels for the two env files being compared.
func Report(w io.Writer, res comparator.Result, leftName, rightName string) {
	if !res.HasDiff() {
		fmt.Fprintln(w, "✔  No differences found.")
		return
	}

	if len(res.MissingInRight) > 0 {
		sort.Strings(res.MissingInRight)
		fmt.Fprintf(w, "Keys in %s but missing in %s:\n", leftName, rightName)
		for _, k := range res.MissingInRight {
			fmt.Fprintf(w, "  - %s\n", k)
		}
	}

	if len(res.MissingInLeft) > 0 {
		sort.Strings(res.MissingInLeft)
		fmt.Fprintf(w, "Keys in %s but missing in %s:\n", rightName, leftName)
		for _, k := range res.MissingInLeft {
			fmt.Fprintf(w, "  + %s\n", k)
		}
	}

	if len(res.Mismatched) > 0 {
		sort.Slice(res.Mismatched, func(i, j int) bool {
			return res.Mismatched[i].Key < res.Mismatched[j].Key
		})
		fmt.Fprintln(w, "Mismatched values:")
		for _, m := range res.Mismatched {
			fmt.Fprintf(w, "  ~ %s\n    %s: %q\n    %s: %q\n",
				m.Key, leftName, m.LeftValue, rightName, m.RightValue)
		}
	}
}
