package main

import (
	"syscall/js"

	"github.com/devfile/library/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	"github.com/feloy/devfile-lifecycle/pkg/graph"
)

func GetFlowChart(this js.Value, args []js.Value) interface{} {
	parserArgs := parser.ParserArgs{
		Data: []byte(args[0].String()),
	}

	devfile, _, err := devfile.ParseDevfileAndValidate(parserArgs)
	if err != nil {
		return ""
	}

	g, err := graph.Build(devfile.Data)
	if err != nil {
		return ""
	}

	return g.ToFlowchart().String()
}

func main() {
	js.Global().Set("getFlowChart", js.FuncOf(GetFlowChart))

	<-make(chan bool)
}
