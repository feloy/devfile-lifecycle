package main

import (
	"syscall/js"

	"github.com/feloy/devfile-lifecycle/pkg/graph"

	apidevfile "github.com/devfile/api/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	context "github.com/devfile/library/v2/pkg/devfile/parser/context"
	"github.com/devfile/library/v2/pkg/testingutil/filesystem"
)

var globalDevfile parser.DevfileObj
var globalFS filesystem.Filesystem

func SetDevfileContent(this js.Value, args []js.Value) interface{} {
	parserArgs := parser.ParserArgs{
		Data: []byte(args[0].String()),
	}

	var err error
	globalDevfile, _, err = devfile.ParseDevfileAndValidate(parserArgs)
	if err != nil {
		return ""
	}
	globalFS = filesystem.NewFakeFs()
	globalDevfile.Ctx = context.FakeContext(globalFS, "/devfile.yaml")

	return getContent()
}

func SetMetadata(this js.Value, args []js.Value) interface{} {
	metadata := args[0]
	globalDevfile.Data.SetMetadata(apidevfile.DevfileMetadata{
		Name:        metadata.Get("name").String(),
		DisplayName: metadata.Get("displayName").String(),
		Description: metadata.Get("description").String(),
	})
	return getContent()
}

func GetFlowChart(this js.Value, args []js.Value) interface{} {
	g, err := graph.Build(globalDevfile.Data)
	if err != nil {
		return ""
	}
	return g.ToFlowchart().String()
}

func getContent() string {
	err := globalDevfile.WriteYamlDevfile()
	if err != nil {
		return ""
	}
	result, err := globalFS.ReadFile("/devfile.yaml")
	if err != nil {
		return ""
	}
	return string(result)
}

func main() {
	js.Global().Set("setDevfileContent", js.FuncOf(SetDevfileContent))
	js.Global().Set("setMetadata", js.FuncOf(SetMetadata))
	js.Global().Set("getFlowChart", js.FuncOf(GetFlowChart))

	<-make(chan bool)
}
