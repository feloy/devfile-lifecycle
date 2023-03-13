package main

import (
	"syscall/js"

	"github.com/devfile/library/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	"github.com/feloy/devfile-lifecycle/pkg/graph"
)

var globalDevfile parser.DevfileObj

func SetDevfileContent(this js.Value, args []js.Value) interface{} {
	parserArgs := parser.ParserArgs{
		Data: []byte(args[0].String()),
	}

	var err error
	globalDevfile, _, err = devfile.ParseDevfileAndValidate(parserArgs)
	if err != nil {
		return ""
	}
	return nil
}

func GetFlowChart(this js.Value, args []js.Value) interface{} {
	g, err := graph.Build(globalDevfile.Data)
	if err != nil {
		return ""
	}
	return g.ToFlowchart().String()
}

func main() {
	js.Global().Set("getFlowChart", js.FuncOf(GetFlowChart))
	js.Global().Set("setDevfileContent", js.FuncOf(SetDevfileContent))

	<-make(chan bool)
}
