package tests

import (
	"bytes"
	"os"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	devfileCtx "github.com/devfile/library/v2/pkg/devfile/parser/context"
	"github.com/devfile/library/v2/pkg/devfile/parser/data"
	"github.com/devfile/library/v2/pkg/testingutil/filesystem"
	"github.com/feloy/devfile-lifecycle/pkg/graph"
	"k8s.io/utils/pointer"
)

func TestBuildGraph(t *testing.T) {
	baseComponent := v1alpha2.Component{
		Name: "my-container",
		ComponentUnion: v1alpha2.ComponentUnion{
			Container: &v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image: "my-image",
				},
				Endpoints: []v1alpha2.Endpoint{
					{
						Name:       "http",
						TargetPort: 8080,
					},
					{
						Name:       "debug",
						TargetPort: 5858,
					},
				},
			},
		},
	}

	secondComponent := v1alpha2.Component{
		Name: "other-container",
		ComponentUnion: v1alpha2.ComponentUnion{
			Container: &v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image: "other-image",
				},
			},
		},
	}

	buildCommand := v1alpha2.Command{
		Id: "my-build",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "go build main.go",
				Component:   "my-container",
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.BuildCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultBuildCommand := *buildCommand.DeepCopy()
	defaultBuildCommand.Exec.Group.IsDefault = pointer.Bool(true)

	buildCommandSecondContainer := v1alpha2.Command{
		Id: "my-build-second-container",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "go build main.go",
				Component:   "other-container",
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.BuildCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultBuildCommandSecondContainer := *buildCommandSecondContainer.DeepCopy()
	defaultBuildCommandSecondContainer.Exec.Group.IsDefault = pointer.Bool(true)
	defaultBuildCommandSecondContainer.Exec.Component = "other-component"

	runCommand := v1alpha2.Command{
		Id: "my-run",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "./main",
				Component:   "my-container",
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.RunCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultRunCommand := *runCommand.DeepCopy()
	defaultRunCommand.Exec.Group.IsDefault = pointer.Bool(true)

	defaultRunHotReloadCommand := *defaultRunCommand.DeepCopy()
	defaultRunHotReloadCommand.Exec.HotReloadCapable = pointer.Bool(true)

	debugCommand := v1alpha2.Command{
		Id: "my-debug",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "./main --debug",
				Component:   "my-container",
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.DebugCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultDebugCommand := *debugCommand.DeepCopy()
	defaultDebugCommand.Exec.Group.IsDefault = pointer.Bool(true)

	build1Command := v1alpha2.Command{
		Id: "my-build-1",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	build2aCommand := v1alpha2.Command{
		Id: "my-build-2a",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 21",
				Component:   "my-container",
			},
		},
	}

	build2bCommand := v1alpha2.Command{
		Id: "my-build-2b",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 22",
				Component:   "my-container",
			},
		},
	}

	build2Command := v1alpha2.Command{
		Id: "my-build-2",
		CommandUnion: v1alpha2.CommandUnion{
			Composite: &v1alpha2.CompositeCommand{
				Commands: []string{
					"my-build-2a",
					"my-build-2b",
				},
			},
		},
	}

	build3Command := v1alpha2.Command{
		Id: "my-build-3",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 3",
				Component:   "my-container",
			},
		},
	}

	compositeBuildCommand := v1alpha2.Command{
		Id: "my-composite-build",
		CommandUnion: v1alpha2.CommandUnion{
			Composite: &v1alpha2.CompositeCommand{
				Commands: []string{
					"my-build-1",
					"my-build-2",
					"my-build-3",
				},
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.BuildCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultCompositeBuildCommand := *compositeBuildCommand.DeepCopy()
	defaultCompositeBuildCommand.Composite.Group.IsDefault = pointer.Bool(true)

	run1Command := v1alpha2.Command{
		Id: "my-run-1",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	run2Command := v1alpha2.Command{
		Id: "my-run-2",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	run3Command := v1alpha2.Command{
		Id: "my-run-3",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	compositeRunCommand := v1alpha2.Command{
		Id: "my-composite-run",
		CommandUnion: v1alpha2.CommandUnion{
			Composite: &v1alpha2.CompositeCommand{
				Commands: []string{
					"my-run-1",
					"my-run-2",
					"my-run-3",
				},
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.RunCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultCompositeRunCommand := *compositeRunCommand.DeepCopy()
	defaultCompositeRunCommand.Composite.Group.IsDefault = pointer.Bool(true)

	debug1Command := v1alpha2.Command{
		Id: "my-debug-1",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	debug2Command := v1alpha2.Command{
		Id: "my-debug-2",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	debug3Command := v1alpha2.Command{
		Id: "my-debug-3",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	compositeDebugCommand := v1alpha2.Command{
		Id: "my-composite-debug",
		CommandUnion: v1alpha2.CommandUnion{
			Composite: &v1alpha2.CompositeCommand{
				Commands: []string{
					"my-debug-1",
					"my-debug-2",
					"my-debug-3",
				},
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.DebugCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultCompositeDebugCommand := *compositeDebugCommand.DeepCopy()
	defaultCompositeDebugCommand.Composite.Group.IsDefault = pointer.Bool(true)

	postStartCommand1 := v1alpha2.Command{
		Id: "post-start-1",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	postStartCommand2 := v1alpha2.Command{
		Id: "post-start-2",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 2",
				Component:   "my-container",
			},
		},
	}

	postStartEvents := []string{"post-start-1", "post-start-2"}

	preStopCommand1 := v1alpha2.Command{
		Id: "pre-stop-1",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 1",
				Component:   "my-container",
			},
		},
	}

	preStopCommand2 := v1alpha2.Command{
		Id: "pre-stop-2",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				CommandLine: "sleep 2",
				Component:   "my-container",
			},
		},
	}

	preStopEvents := []string{"pre-stop-1", "pre-stop-2"}

	deployCommand := v1alpha2.Command{
		Id: "my-deploy",
		CommandUnion: v1alpha2.CommandUnion{
			Apply: &v1alpha2.ApplyCommand{
				Component: "kubernetes-deploy",
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.DeployCommandGroupKind,
						},
					},
				},
			},
		},
	}

	defaultDeployCommand := *deployCommand.DeepCopy()
	defaultDeployCommand.Apply.Group.IsDefault = pointer.Bool(true)

	tests := []struct {
		name            string
		filename        string
		dataVersion     string
		components      func() []v1alpha2.Component
		commands        func() []v1alpha2.Command
		postStartEvents []string
		preStopEvents   []string
	}{
		{
			name:        "container only",
			filename:    "container-only",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands:    func() []v1alpha2.Command { return nil },
		},
		{
			name:        "container with Exec Build and Run commands",
			filename:    "container-build-run",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunCommand,
				}
			},
		},
		{
			name:        "container with Exec Build, Run and Debug commands",
			filename:    "container-build-run-debug",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunCommand,
					defaultDebugCommand,
				}
			},
		},
		{
			name:        "container with Exec Build, Run and Debug commands with 2 containers",
			filename:    "container-build-run-debug-2-containers",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent, secondComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommandSecondContainer,
					defaultRunCommand,
					defaultDebugCommand,
				}
			},
		},
		{
			name:        "container with Exec Build and HotReload Capable Run commands",
			filename:    "container-build-run-hot-reload",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunHotReloadCommand,
				}
			},
		},
		{
			name:        "container with Composite Build and Exec Run commands",
			filename:    "container-composite-build-exec-run",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					build2aCommand,
					build2bCommand,
					build1Command,
					build2Command,
					build3Command,
					defaultCompositeBuildCommand,
					defaultRunCommand,
				}
			},
		},
		{
			name:        "container with Composite Build and Run commands",
			filename:    "container-composite-build-run",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					build2aCommand,
					build2bCommand,
					build1Command,
					build2Command,
					build3Command,
					run1Command,
					run2Command,
					run3Command,
					defaultCompositeBuildCommand,
					defaultCompositeRunCommand,
				}
			},
		},
		{
			name:        "container with Composite Build, Run and Debug commands",
			filename:    "container-composite-build-run-debug",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					build2aCommand,
					build2bCommand,
					build1Command,
					build2Command,
					build3Command,
					defaultCompositeBuildCommand,
					run1Command,
					run2Command,
					run3Command,
					defaultCompositeRunCommand,
					debug1Command,
					debug2Command,
					debug3Command,
					defaultCompositeDebugCommand,
				}
			},
		},
		{
			name:        "container with Exec Build and Run commands, and postStart event",
			filename:    "container-build-run-post-start",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunCommand,
					postStartCommand1,
					postStartCommand2,
				}
			},
			postStartEvents: postStartEvents,
		},
		{
			name:        "container with Exec Build and Run commands, and preStop event",
			filename:    "container-build-run-pre-stop",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunCommand,
					preStopCommand1,
					preStopCommand2,
				}
			},
			preStopEvents: preStopEvents,
		},
		{
			name:        "container with Exec Build and Run commands + a deploy command",
			filename:    "container-build-run-deploy",
			dataVersion: string(data.APISchemaVersion200),
			components:  func() []v1alpha2.Component { return []v1alpha2.Component{baseComponent} },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunCommand,
					defaultDeployCommand,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			devfileData, err := data.NewDevfileData(tt.dataVersion)
			if err != nil {
				t.Error(err)
			}
			devfileData.SetSchemaVersion(tt.dataVersion)

			err = devfileData.AddComponents(tt.components())
			if err != nil {
				t.Error(err)
			}

			err = devfileData.AddCommands(tt.commands())
			if err != nil {
				t.Error(err)
			}

			if tt.postStartEvents != nil || tt.preStopEvents != nil {
				events := v1alpha2.Events{
					DevWorkspaceEvents: v1alpha2.DevWorkspaceEvents{
						PostStart: tt.postStartEvents,
						PreStop:   tt.preStopEvents,
					},
				}

				err = devfileData.AddEvents(events)
				if err != nil {
					t.Error(err)
				}
			}

			devfileBytes := getDevfileBytes(t, devfileData)
			shotBytes := getShotBytes(tt.filename, "devfiles/", ".yaml")

			if !bytes.Equal(devfileBytes, shotBytes) {
				t.Errorf("Devfile should be:\n%s\nbut is:\n%s\n", devfileBytes, shotBytes)
			}

			graphBytes := getGraphBytes(t, devfileData)
			graphShotBytes := getShotBytes(tt.filename, "graphs/", ".md")
			if !bytes.Equal(graphBytes, graphShotBytes) {
				t.Errorf("Graph should be:\n%s\nbut is:\n%s\n", graphBytes, graphShotBytes)
			}
		})
	}
}

func getDevfileBytes(t *testing.T, devfileData data.DevfileData) []byte {
	fs := filesystem.NewFakeFs()
	filename := "/devfile.yaml"
	devfileObj := parser.DevfileObj{
		Data: devfileData,
		Ctx:  devfileCtx.FakeContext(fs, filename),
	}
	err := devfileObj.WriteYamlDevfile()
	if err != nil {
		t.Error(err)
	}

	devfileContent, err := fs.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	return devfileContent
}

func getShotBytes(filename, prefix, suffix string) []byte {
	content, err := os.ReadFile(prefix + filename + suffix)
	if err != nil {
		return []byte{}
	}
	return content
}

func getGraphBytes(t *testing.T, devfileData data.DevfileData) []byte {
	g, err := graph.Build(devfileData)
	if err != nil {
		t.Error(err)
	}
	return []byte("```mermaid\n" + g.ToFlowchart() + "```\n")
}
