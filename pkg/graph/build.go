package graph

import (
	"errors"
	"fmt"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser/data"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	"k8s.io/utils/pointer"

	"github.com/feloy/devfile-lifecycle/pkg/dftools"
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

	start := g.AddNode("start")

	containerNode := g.AddNode(
		container.Name,
		"container: "+container.Name,
		"image: "+container.Container.Image,
	)
	g.EntryNodeID = containerNode.ID

	_ = g.AddEdge(
		start,
		containerNode,
		"dev",
	)

	syncNodeStart := g.AddNode(
		"sync-all-"+container.Name,
		"Sync All Sources",
	)

	_ = g.AddEdge(
		containerNode,
		syncNodeStart,
		"container running",
	)

	/* Get PostStart event */
	postStartEvents := devfileData.GetEvents().PostStart

	previousNode := syncNodeStart
	nextText := "sources synced"
	for _, postStartEvent := range postStartEvents {
		node := g.AddNode(postStartEvent, "Post Start", "command: "+postStartEvent)
		_ = g.AddEdge(
			previousNode,
			node,
			nextText,
		)
		previousNode = node
		nextText = postStartEvent + " done"
	}

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

	buildNodesStart, buildNodesEnd, err := addCommand(g, devfileData, defaultBuildCommand, []*Node{previousNode}, nextText)
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

		runNodesStart, runNodesEnd, err := addCommand(g, devfileData, defaultRunCommand, buildNodesEnd, edgeText)
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
			container.Name+"-"+runNodesStart[0].ID+"-expose",
			lines...,
		)

		for _, runNodeEnd := range runNodesEnd {
			_ = g.AddEdge(
				runNodeEnd,
				exposeNode,
				"command running",
			)
		}

		/* Get PreStop event */
		preStopEvents := devfileData.GetEvents().PreStop

		previousNode := exposeNode
		nextText := "User quits"
		for _, preStopEvent := range preStopEvents {
			node := g.AddNode(preStopEvent, "Pre Stop", "command: "+preStopEvent)
			_ = g.AddEdge(
				previousNode,
				node,
				nextText,
			)
			previousNode = node
			nextText = preStopEvent + " done"
		}

		/* Add "stop container" node */

		stopNode := g.AddNode(container.Name+"-stop", "Stop container", "container: "+container.Name)

		/* Add "user quits" edge */

		_ = g.AddEdge(
			previousNode,
			stopNode,
			nextText,
		)

		_, syncNodeChangedExists := g.nodes["sync-modified-"+container.Name]

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

		if !syncNodeChangedExists {
			if defaultRunCommand.Exec != nil && defaultRunCommand.Exec.HotReloadCapable != nil && *defaultRunCommand.Exec.HotReloadCapable {
				_ = g.AddEdge(
					syncNodeChanged,
					exposeNode,
					"source synced",
				)
			} else {
				for _, buildNodeStart := range buildNodesStart {
					_ = g.AddEdge(
						syncNodeChanged,
						buildNodeStart,
						"source synced",
					)
				}
			}
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

func addCommand(g *Graph, devfileData data.DevfileData, command v1alpha2.Command, nodesBefore []*Node, text ...string) (start []*Node, end []*Node, err error) {
	if command.Exec != nil {
		return addExecCommand(g, command, nodesBefore, text...)
	}
	if command.Composite != nil {
		return addCompositeCommand(g, devfileData, command, nodesBefore, text...)
	}
	return nil, nil, fmt.Errorf("Command type not implemented for %s", command.Id)
}

func addExecCommand(g *Graph, command v1alpha2.Command, nodesBefore []*Node, text ...string) ([]*Node, []*Node, error) {
	node := g.AddNode(
		command.Id,
		"command: "+command.Id,
	)

	for _, nodeBefore := range nodesBefore {
		_ = g.AddEdge(
			nodeBefore,
			node,
			text...,
		)
	}

	return []*Node{node}, []*Node{node}, nil

}

func addCompositeCommand(g *Graph, devfileData data.DevfileData, command v1alpha2.Command, nodesBefore []*Node, text ...string) ([]*Node, []*Node, error) {
	// Serial
	if !pointer.BoolDeref(command.Composite.Parallel, false) {
		previousNodes := nodesBefore
		var firstNode []*Node
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
			var first []*Node
			first, previousNodes, err = addCommand(g, devfileData, subcommands[0], previousNodes, text...)
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

		return firstNode, previousNodes, nil
	}

	// Parallel

	startNodes := make([]*Node, 0, len(command.Composite.Commands))
	endNodes := make([]*Node, 0, len(command.Composite.Commands))

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
		start, end, err := addCommand(g, devfileData, subcommands[0], nodesBefore, text...)
		if err != nil {
			return nil, nil, err
		}
		startNodes = append(startNodes, start...)
		endNodes = append(endNodes, end...)
	}
	return startNodes, endNodes, nil
}
