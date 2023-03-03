package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/devfile/library/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser"

	"github.com/feloy/devfile-lifecycle/pkg/graph"
)

func usage(args []string, withError bool) {
	fmt.Fprintf(os.Stderr, "Usage: %s devfile.yaml", args[0])
	if withError {
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) != 2 {
		usage(os.Args, true)
	}
	file := os.Args[1]
	parserArgs := parser.ParserArgs{}
	if strings.HasPrefix(file, "http") {
		parserArgs.URL = file
	} else {
		parserArgs.Path = file
	}
	devfile, _, err := devfile.ParseDevfileAndValidate(parserArgs)
	if err != nil {
		panic(err)
	}

	g, err := graph.Build(devfile.Data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("```mermaid\n%s```\n", g.ToFlowchart())
}
