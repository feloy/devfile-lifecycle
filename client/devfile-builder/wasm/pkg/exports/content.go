package exports

import (
	"errors"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apidevfile "github.com/devfile/api/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	"k8s.io/utils/pointer"

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
		commands := make([]interface{}, len(container.ComponentUnion.Container.Command))
		for i, command := range container.ComponentUnion.Container.Command {
			commands[i] = command
		}

		args := make([]interface{}, len(container.ComponentUnion.Container.Args))
		for i, arg := range container.ComponentUnion.Container.Args {
			args[i] = arg
		}

		userCommands, err := getUserCommands(container.Name)
		if err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"name":         container.Name,
			"image":        container.ComponentUnion.Container.Image,
			"command":      commands,
			"args":         args,
			"userCommands": userCommands,
		})
	}
	return result, nil
}

func getUserCommands(component string) ([]interface{}, error) {
	result := []interface{}{}

	commands, err := global.Devfile.Data.GetCommands(common.DevfileOptions{
		CommandOptions: common.CommandOptions{
			CommandType: v1alpha2.ExecCommandType,
		},
	})
	if err != nil {
		return nil, err
	}
	for _, command := range commands {
		if command.Exec.Component != component {
			continue
		}
		result = append(result, map[string]interface{}{
			"name":             command.Id,
			"commandLine":      command.CommandUnion.Exec.CommandLine,
			"hotReloadCapable": pointer.BoolDeref(command.CommandUnion.Exec.HotReloadCapable, false),
			"workingDir":       command.CommandUnion.Exec.WorkingDir,
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
