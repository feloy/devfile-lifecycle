package graph

import (
	"bytes"
	"os"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	devfileCtx "github.com/devfile/library/v2/pkg/devfile/parser/context"
	"github.com/devfile/library/v2/pkg/devfile/parser/data"
	"github.com/devfile/library/v2/pkg/testingutil/filesystem"
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

	tests := []struct {
		name        string
		filename    string
		dataVersion string
		component   func() v1alpha2.Component
		commands    func() []v1alpha2.Command
	}{
		{
			name:        "container only",
			filename:    "container-only",
			dataVersion: string(data.APISchemaVersion200),
			component:   func() v1alpha2.Component { return baseComponent },
			commands:    func() []v1alpha2.Command { return nil },
		},
		{
			name:        "container with Exec Build and Run commands",
			filename:    "container-build-run",
			dataVersion: string(data.APISchemaVersion200),
			component:   func() v1alpha2.Component { return baseComponent },
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
			component:   func() v1alpha2.Component { return baseComponent },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunCommand,
					defaultDebugCommand,
				}
			},
		},
		{
			name:        "container with Exec Build and HotReload Capable Run commands",
			filename:    "container-build-run-hot-reload",
			dataVersion: string(data.APISchemaVersion200),
			component:   func() v1alpha2.Component { return baseComponent },
			commands: func() []v1alpha2.Command {
				return []v1alpha2.Command{
					defaultBuildCommand,
					defaultRunHotReloadCommand,
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

			err = devfileData.AddComponents([]v1alpha2.Component{tt.component()})
			if err != nil {
				t.Error(err)
			}

			err = devfileData.AddCommands(tt.commands())
			if err != nil {
				t.Error(err)
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

	//err = devfileData.AddCommands(tt.execCommands)
	//if err != nil {
	//	t.Error(err)
	//}
	//err = devfileData.AddCommands(tt.compCommands)
	//if err != nil {
	//	t.Error(err)
	//}

	// fmt.Printf("%s\n", s)
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
	g, err := Build(devfileData)
	if err != nil {
		t.Error(err)
	}
	return []byte("```mermaid\n" + g.ToFlowchart().String() + "```\n")
}
