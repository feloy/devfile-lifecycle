package main

import (
	"errors"
	"strings"
	"syscall/js"

	"github.com/feloy/devfile-lifecycle/pkg/graph"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apidevfile "github.com/devfile/api/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	context "github.com/devfile/library/v2/pkg/devfile/parser/context"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	"github.com/devfile/library/v2/pkg/testingutil/filesystem"
)

const (
	SEPARATOR = ","
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
		Name:              metadata.Get("name").String(),
		Version:           metadata.Get("version").String(),
		DisplayName:       metadata.Get("displayName").String(),
		Description:       metadata.Get("description").String(),
		Tags:              splitTags(metadata.Get("tags").String()),
		Architectures:     splitArchitectures(metadata.Get("architectures").String()),
		Icon:              metadata.Get("icon").String(),
		GlobalMemoryLimit: metadata.Get("globalMemoryLimit").String(),
		ProjectType:       metadata.Get("projectType").String(),
		Language:          metadata.Get("language").String(),
		Website:           metadata.Get("website").String(),
		Provider:          metadata.Get("provider").String(),
		SupportUrl:        metadata.Get("supportUrl").String(),
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

	devEnvs, err := getDevEnvs()
	if err != nil {
		return nil, errors.New("error getting development environments")
	}

	return map[string]interface{}{
		"content":  string(result),
		"metadata": getMetadata(),
		"devEnvs":  devEnvs,
	}, nil
}

func getMetadata() map[string]interface{} {
	metadata := globalDevfile.Data.GetMetadata()
	return map[string]interface{}{
		"name":              metadata.Name,
		"version":           metadata.Version,
		"displayName":       metadata.DisplayName,
		"description":       metadata.Description,
		"tags":              strings.Join(metadata.Tags, SEPARATOR),
		"architectures":     joinArchitectures(metadata.Architectures),
		"icon":              metadata.Icon,
		"globalMemoryLimit": metadata.GlobalMemoryLimit,
		"projectType":       metadata.ProjectType,
		"language":          metadata.Language,
		"website":           metadata.Website,
		"provider":          metadata.Provider,
		"supportUrl":        metadata.SupportUrl,
	}
}

func getDevEnvs() ([]interface{}, error) {
	containers, err := globalDevfile.Data.GetComponents(common.DevfileOptions{
		ComponentOptions: common.ComponentOptions{
			ComponentType: v1alpha2.ContainerComponentType,
		},
	})
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, 0, len(containers))
	for _, container := range containers {
		result = append(result, map[string]interface{}{
			"name":  container.Name,
			"image": container.ComponentUnion.Container.Image,
		})
	}
	return result, nil
}

func joinArchitectures(architectures []apidevfile.Architecture) string {
	strArchs := make([]string, len(architectures))
	for i, arch := range architectures {
		strArchs[i] = string(arch)
	}
	return strings.Join(strArchs, SEPARATOR)
}

func splitArchitectures(architectures string) []apidevfile.Architecture {
	if architectures == "" {
		return nil
	}
	parts := strings.Split(architectures, SEPARATOR)
	result := make([]apidevfile.Architecture, len(parts))
	for i, arch := range parts {
		result[i] = apidevfile.Architecture(strings.Trim(arch, " "))
	}
	return result
}

func splitTags(tags string) []string {
	if tags == "" {
		return nil
	}
	parts := strings.Split(tags, SEPARATOR)
	result := make([]string, len(parts))
	for i, tag := range parts {
		result[i] = strings.Trim(tag, " ")
	}
	return result
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
