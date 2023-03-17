package exports

import (
	"syscall/js"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/feloy/devfile-lifecycle/client/devfile-builder/wasm/pkg/global"
)

func UpdateContainerWrapper(this js.Value, args []js.Value) interface{} {
	command := getStringArray(args[2])
	arg := getStringArray(args[3])
	return result(
		updateContainer(args[0].String(), args[1].String(), command, arg),
	)
}

func updateContainer(name string, image string, command []string, args []string) (map[string]interface{}, error) {
	component := v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Container: &v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image:   image,
					Command: command,
					Args:    args,
				},
			},
		},
	}
	global.Devfile.Data.UpdateComponent(component)
	return getContent()
}
