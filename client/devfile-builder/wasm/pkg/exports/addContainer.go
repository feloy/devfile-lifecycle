package exports

import (
	"syscall/js"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/feloy/devfile-lifecycle/client/devfile-builder/wasm/pkg/global"
)

func AddContainerWrapper(this js.Value, args []js.Value) interface{} {
	return result(
		addContainer(args[0].String(), args[1].String()),
	)
}

func addContainer(name string, image string) (map[string]interface{}, error) {
	container := v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Container: &v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image: image,
				},
			},
		},
	}
	err := global.Devfile.Data.AddComponents([]v1alpha2.Component{container})
	if err != nil {
		return nil, err
	}
	return getContent()
}
