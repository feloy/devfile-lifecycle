package main

import (
	"errors"
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

// setDevfileContent

func SetDevfileContentWrapper(this js.Value, args []js.Value) interface{} {
	return result(
		setDevfileContent(args[0].String()),
	)
}

func setDevfileContent(content string) (map[string]interface{}, error) {
	parserArgs := parser.ParserArgs{
		Data: []byte(content),
	}
	var err error
	globalDevfile, _, err = devfile.ParseDevfileAndValidate(parserArgs)
	if err != nil {
		return nil, errors.New("error parsing devfile")
	}
	globalFS = filesystem.NewFakeFs()
	globalDevfile.Ctx = context.FakeContext(globalFS, "/devfile.yaml")

	return getContent()
}

// setMetadata

func SetMetadataWrapper(this js.Value, args []js.Value) interface{} {
	return result(
		setMetadata(args[0]),
	)
}

func setMetadata(metadata js.Value) (map[string]interface{}, error) {
	globalDevfile.Data.SetMetadata(apidevfile.DevfileMetadata{
		Name:        metadata.Get("name").String(),
		DisplayName: metadata.Get("displayName").String(),
		Description: metadata.Get("description").String(),
	})
	return getContent()
}

// getFlowChart

func GetFlowChartWrapper(this js.Value, args []js.Value) interface{} {
	return result(
		getFlowChart(),
	)
}

func getFlowChart() (string, error) {
	g, err := graph.Build(globalDevfile.Data)
	if err != nil {
		return "", errors.New("error building graph")
	}
	return g.ToFlowchart().String(), nil
}

// common

// getContent returns the YAML content of the global devfile as string
func getContent() (map[string]interface{}, error) {
	err := globalDevfile.WriteYamlDevfile()
	if err != nil {
		return nil, errors.New("error writing file")
	}
	result, err := globalFS.ReadFile("/devfile.yaml")
	if err != nil {
		return nil, errors.New("error reading file")
	}

	metadata := globalDevfile.Data.GetMetadata()
	metadataResult := map[string]interface{}{
		"name":        metadata.Name,
		"displayName": metadata.DisplayName,
		"description": metadata.Description,
	}
	return map[string]interface{}{
		"content":  string(result),
		"metadata": metadataResult,
	}, nil
}

// result returns the value and error in a format acceptable for JS
func result(value interface{}, err error) map[string]interface{} {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	return map[string]interface{}{
		"value": value,
		"err":   errStr,
	}
}

func main() {
	js.Global().Set("setDevfileContent", js.FuncOf(SetDevfileContentWrapper))
	js.Global().Set("setMetadata", js.FuncOf(SetMetadataWrapper))
	js.Global().Set("getFlowChart", js.FuncOf(GetFlowChartWrapper))

	<-make(chan bool)
}
