package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gostaticanalysis/knife/lib/knife"
)

var (
	flagFormat string
	flagXPath  string
)

func init() {
	flag.StringVar(&flagFormat, "f", "{{.}}", "output format")
	flag.StringVar(&flagXPath, "xpath", "", "A XPath expression for an AST node")
	flag.Parse()
}

func main() {
	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	k, err := knife.New(args...)
	if err != nil {
		return err
	}

	if flagXPath != "" {
		err := k.FormatWithXPath(os.Stdout, flagFormat, flagXPath)
		if err != nil {
			return err
		}
		return nil
	}

	return k.Format(os.Stdout, flagFormat)
}
