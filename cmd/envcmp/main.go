package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/filter"
	"github.com/user/envcmp/internal/formatter"
	"github.com/user/envcmp/internal/parser"
)

func main() {
	strict := flag.Bool("strict", false, "exit 1 if any diff is found")
	onlyMissing := flag.Bool("only-missing", false, "report only missing keys")
	onlyMismatched := flag.Bool("only-mismatched", false, "report only mismatched values")
	keyFilter := flag.String("keys", "", "comma-separated list of keys to include")
	format := flag.String("format", "text", "output format: text, json, markdown")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: envcmp [flags] <file1> <file2>")
		os.Exit(2)
	}

	left, err := parser.ParseFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[0], err)
		os.Exit(2)
	}

	right, err := parser.ParseFile(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", args[1], err)
		os.Exit(2)
	}

	result := comparator.Compare(left, right)

	result = filter.Apply(result, filter.Options{
		OnlyMissing:    *onlyMissing,
		OnlyMismatched: *onlyMismatched,
		Keys:           *keyFilter,
	})

	out, err := formatter.Render(result, formatter.Format(*format))
	if err != nil {
		fmt.Fprintf(os.Stderr, "formatter error: %v\n", err)
		os.Exit(2)
	}

	if out != "" {
		fmt.Print(out)
	}

	if *strict && !filter.IsEmpty(result) {
		os.Exit(1)
	}
}
