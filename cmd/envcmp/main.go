package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envcmp/internal/comparator"
	"github.com/user/envcmp/internal/parser"
	"github.com/user/envcmp/internal/reporter"
)

const usage = `envcmp - diff .env files across environments

Usage:
  envcmp <file1> <file2> [flags]

Flags:
`

func main() {
	quiet := flag.Bool("quiet", false, "suppress output, exit 1 if differences found")
	strict := flag.Bool("strict", false, "exit 1 if any differences are found")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	leftPath := args[0]
	rightPath := args[1]

	left, err := parser.ParseFile(leftPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", leftPath, err)
		os.Exit(2)
	}

	right, err := parser.ParseFile(rightPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", rightPath, err)
		os.Exit(2)
	}

	result := comparator.Compare(left, right)

	if !*quiet {
		reporter.Report(os.Stdout, leftPath, rightPath, result)
	}

	if *strict && result.HasDiff() {
		os.Exit(1)
	}
}
