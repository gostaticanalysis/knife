package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gostaticanalysis/knife"
	"github.com/gostaticanalysis/knife/mcp"
)

var (
	flagVersion   bool
	flagFormat    string
	flagTemplate  string
	flagExtraData string
	flagXPath     string
	flagTests     bool
)

func init() {
	flag.BoolVar(&flagVersion, "v", false, "print version")
	flag.StringVar(&flagFormat, "f", "{{.}}", "output format")
	flag.StringVar(&flagTemplate, "template", "", "template file")
	flag.StringVar(&flagExtraData, "data", "", "extra data (key:value,key:value)")
	flag.StringVar(&flagXPath, "xpath", "", "A XPath expression for an AST node")
	flag.BoolVar(&flagTests, "tests", true, "include test files")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	if err := run(ctx, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if flagVersion {
		fmt.Println("knife", knife.Version())
		return nil
	}

	if len(args) > 0 && args[0] == "mcp" {
		return runMCPServer(ctx)
	}

	knifeOpt := &knife.KnifeOption{
		Tests: flagTests,
	}
	k, err := knife.New(knifeOpt, args...)
	if err != nil {
		return err
	}

	var w io.Writer = os.Stdout

	opt := &knife.ExecuteOption{
		XPath: flagXPath,
	}

	if flagExtraData != "" {
		extraData, err := parseExtraData(flagExtraData)
		if err != nil {
			return err
		}
		opt.ExtraData = extraData
	}

	var tmpl any = flagFormat
	if flagTemplate != "" {
		tmpl, err = os.ReadFile(flagTemplate)
		if err != nil {
			return fmt.Errorf("cannot read template: %w", err)
		}
	}

	pkgs := k.Packages()
	for i, pkg := range pkgs {

		if err := k.Execute(w, pkg, tmpl, opt); err != nil {
			return err
		}

		if i != len(pkgs)-1 {
			fmt.Fprintln(w)
		}
	}

	return nil
}

func parseExtraData(extraData string) (map[string]any, error) {
	m := map[string]any{}
	kvs := strings.Split(extraData, ",")
	for i := range kvs {
		kv := strings.Split(strings.TrimSpace(kvs[i]), ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid extra data: %s", kvs[i])
		}
		m[kv[0]] = kv[1]
	}
	return m, nil
}

func runMCPServer(ctx context.Context) error {
	server := mcp.NewKnifeServer()
	return server.Run(ctx, mcpsdk.NewStdioTransport())
}
