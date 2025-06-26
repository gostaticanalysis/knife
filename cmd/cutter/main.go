package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/newmo-oss/gogroup"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/gostaticanalysis/knife"
	"github.com/gostaticanalysis/knife/cutter"
	"github.com/gostaticanalysis/knife/mcp"
)

var (
	flagVersion   bool
	flagFormat    string
	flagTemplate  string
	flagExtraData string
)

func init() {
	flag.BoolVar(&flagVersion, "v", false, "print version")
	flag.StringVar(&flagFormat, "f", "{{.}}", "output format")
	flag.StringVar(&flagTemplate, "template", "", "template file")
	flag.StringVar(&flagExtraData, "data", "", "extra data (key:value,key:value)")
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
		fmt.Println("cutter", knife.Version())
		return nil
	}

	if len(args) > 0 && args[0] == "mcp" {
		return runMCPServer(ctx)
	}

	c, err := cutter.New(args...)
	if err != nil {
		return err
	}

	var opt cutter.Option
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

	pkgs := c.KnifePackages()
	readers := make([]io.Reader, len(pkgs))
	var g gogroup.Group
	for i, pkg := range pkgs {
		g.Add(func(ctx context.Context) error {
			var buf bytes.Buffer
			if err := c.Execute(&buf, pkg, tmpl, &opt); err != nil {
				return err
			}

			if i != len(pkgs)-1 {
				if _, err := fmt.Fprintln(&buf); err != nil {
					return err
				}
			}

			// no race condition
			readers[i] = &buf

			return nil
		})
	}

	if err := g.Run(ctx); err != nil {
		return err
	}

	if _, err := io.Copy(os.Stdout, io.MultiReader(readers...)); err != nil {
		return err
	}

	return nil
}

func parseExtraData(extraData string) (map[string]any, error) {
	m := make(map[string]any)
	kvs := strings.Split(extraData, ",")
	for i := range kvs {
		key, value, ok := strings.Cut(strings.TrimSpace(kvs[i]), ":")
		if !ok {
			return nil, fmt.Errorf("invalid extra data: %s", kvs[i])
		}
		m[key] = value
	}
	return m, nil
}

func runMCPServer(ctx context.Context) error {
	server := mcp.NewCutterServer()
	return server.Run(ctx, mcpsdk.NewStdioTransport())
}
