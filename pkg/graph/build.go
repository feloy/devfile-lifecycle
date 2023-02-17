package graph

import (
	"errors"
	"fmt"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser/data"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"

	"github.com/feloy/devfile-states/pkg/dftools"
)

func Build(devfileData data.DevfileData) (*Graph, error) {

	g := NewGraph()

	/* Get "container" node */
	containers, err := devfileData.GetComponents(common.DevfileOptions{
		ComponentOptions: common.ComponentOptions{
			ComponentType: v1alpha2.ContainerComponentType,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(containers) != 1 {
		return nil, errors.New("more than one container, not supported yet")
	}

	container := containers[0]

	containerNode := g.AddNode(
		container.Name,
		"container: "+container.Name,
		"image: "+container.Container.Image,
	)
	g.EntryNodeID = containerNode.ID

	/* Get "build command" node */
	buildCommands, err := devfileData.GetCommands(common.DevfileOptions{
		CommandOptions: common.CommandOptions{
			CommandGroupKind: v1alpha2.BuildCommandGroupKind,
		},
	})
	if err != nil {
		return nil, err
	}

	var defaultBuildCommand v1alpha2.Command
	for _, buildCommand := range buildCommands {
		if dftools.GetCommandGroup(buildCommand).IsDefault != nil && *dftools.GetCommandGroup(buildCommand).IsDefault {
			defaultBuildCommand = buildCommand
			break
		}
	}

	if defaultBuildCommand.Id == "" {
		return g, nil
	}

	buildNode := g.AddNode(
		defaultBuildCommand.Id,
		"command: "+defaultBuildCommand.Id,
	)

	_ = g.AddEdge(
		containerNode,
		buildNode,
		"container running",
	)

	for _, debug := range []bool{false, true} {
		/* Get "run command" node */

		kind := v1alpha2.RunCommandGroupKind
		if debug {
			kind = v1alpha2.DebugCommandGroupKind
		}
		runCommands, err := devfileData.GetCommands(common.DevfileOptions{
			CommandOptions: common.CommandOptions{
				CommandGroupKind: kind,
			},
		})
		if err != nil {
			return nil, err
		}

		var defaultRunCommand v1alpha2.Command
		for _, runCommand := range runCommands {
			if dftools.GetCommandGroup(runCommand).IsDefault != nil && *dftools.GetCommandGroup(runCommand).IsDefault {
				defaultRunCommand = runCommand
				break
			}
		}

		if defaultRunCommand.Id == "" {
			continue
		}

		runNode := g.AddNode(
			defaultRunCommand.Id,
			"command: "+defaultRunCommand.Id,
		)

		edgeText := "build done, "
		if debug {
			edgeText += "with debug"
		} else {
			edgeText += "with run"
		}
		_ = g.AddEdge(
			buildNode,
			runNode,
			edgeText,
		)

		lines := []string{
			"Expose ports",
		}
		for _, endpoint := range container.Container.Endpoints {
			if !debug && dftools.IsDebugPort(endpoint) {
				continue
			}
			lines = append(lines, fmt.Sprintf("%s: %d", endpoint.Name, endpoint.TargetPort))
		}
		exposeNode := g.AddNode(
			container.Name+"-"+runNode.ID+"-expose",
			lines...,
		)

		_ = g.AddEdge(
			runNode,
			exposeNode,
			"command running",
		)

		/* Add "source changed" edge */

		_ = g.AddEdge(
			exposeNode,
			buildNode,
			"source changed",
		)

		/* Add "devfile changed" edge */

		_ = g.AddEdge(
			exposeNode,
			containerNode,
			"devfile changed",
		)
	}

	return g, nil
}
