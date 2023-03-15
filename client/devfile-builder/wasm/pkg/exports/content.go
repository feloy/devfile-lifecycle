package exports

import (
	"errors"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apidevfile "github.com/devfile/api/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"

	"github.com/feloy/devfile-lifecycle/client/devfile-builder/wasm/pkg/global"
)

const (
	SEPARATOR = ","
)

// getContent returns the YAML content of the global devfile as string
func getContent() (map[string]interface{}, error) {
	err := global.Devfile.WriteYamlDevfile()
	if err != nil {
		return nil, errors.New("error writing file")
	}
	result, err := global.FS.ReadFile("/devfile.yaml")
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
	metadata := global.Devfile.Data.GetMetadata()
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
	containers, err := global.Devfile.Data.GetComponents(common.DevfileOptions{
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
