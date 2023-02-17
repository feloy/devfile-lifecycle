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

	syncNodeStart := g.AddNode(
		"sync-all-"+container.Name,
		"Sync All Sources",
	)

	_ = g.AddEdge(
		containerNode,
		syncNodeStart,
		"container running",
	)

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

	buildNodeStart, buildNodeEnd, err := addCommand(g, devfileData, defaultBuildCommand, syncNodeStart, "sources synced")
	if err != nil {
		return nil, err
	}

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

		edgeText := "build done, "
		if debug {
			edgeText += "with debug"
		} else {
			edgeText += "with run"
		}

		runNode, _, err := addCommand(g, devfileData, defaultRunCommand, buildNodeEnd, edgeText)
		if err != nil {
			return nil, err
		}

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

		// Add "Sync" node
		syncNodeChanged := g.AddNode(
			"sync-modified-"+container.Name,
			"Sync Modified Sources",
		)

		_ = g.AddEdge(
			exposeNode,
			syncNodeChanged,
			"source changed",
		)

		/* Add "source synced" edge */

		if defaultRunCommand.Exec.HotReloadCapable != nil && *defaultRunCommand.Exec.HotReloadCapable {
			_ = g.AddEdge(
				syncNodeChanged,
				exposeNode,
				"source synced",
			)
		} else {
			_ = g.AddEdge(
				syncNodeChanged,
				buildNodeStart,
				"source synced",
			)
		}

		/* Add "devfile changed" edge */

		_ = g.AddEdge(
			exposeNode,
			containerNode,
			"devfile changed",
		)
	}

	return g, nil
}

func addCommand(g *Graph, devfileData data.DevfileData, command v1alpha2.Command, nodeBefore *Node, text ...string) (start *Node, end *Node, err error) {
	if command.Exec != nil {
		return addExecCommand(g, command, nodeBefore, text...)
	}
	if command.Composite != nil {
		return addCompositeCommand(g, devfileData, command, nodeBefore, text...)
	}
	return nil, nil, fmt.Errorf("Command type not implemented for %s", command.Id)
}

func addExecCommand(g *Graph, command v1alpha2.Command, nodeBefore *Node, text ...string) (*Node, *Node, error) {
	node := g.AddNode(
		command.Id,
		"command: "+command.Id,
	)

	_ = g.AddEdge(
		nodeBefore,
		node,
		text...,
	)

	return node, node, nil

}

func addCompositeCommand(g *Graph, devfileData data.DevfileData, command v1alpha2.Command, nodeBefore *Node, text ...string) (*Node, *Node, error) {
	previousNode := nodeBefore
	var firstNode *Node
	for _, subcommandName := range command.Composite.Commands {
		subcommands, err := devfileData.GetCommands(common.DevfileOptions{
			FilterByName: subcommandName,
		})
		if err != nil {
			return nil, nil, err
		}
		if len(subcommands) != 1 {
			return nil, nil, fmt.Errorf("command not found: %s", subcommandName)
		}
		var first *Node
		first, previousNode, err = addCommand(g, devfileData, subcommands[0], previousNode, text...)
		if err != nil {
			return nil, nil, err
		}
		if firstNode == nil {
			firstNode = first
		}
		text = []string{
			subcommandName + " done",
		}
	}

	return firstNode, previousNode, nil
}
