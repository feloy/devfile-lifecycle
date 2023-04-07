package main

import (
	"syscall/js"

	"github.com/feloy/devfile-lifecycle/client/devfile-builder/wasm/pkg/exports"
)

func setFreshDevfile() {
	content := `schemaVersion: 2.2.0`
	exports.SetDevfileContent(content)
}

func main() {
	setFreshDevfile()

	js.Global().Set("addContainer", js.FuncOf(exports.AddContainerWrapper))
	js.Global().Set("addExecCommand", js.FuncOf(exports.AddExecCommandWrapper))
	js.Global().Set("addApplyCommand", js.FuncOf(exports.AddApplyCommandWrapper))
	js.Global().Set("addUserCommand", js.FuncOf(exports.AddUserCommandWrapper))
	js.Global().Set("getFlowChart", js.FuncOf(exports.GetFlowChartWrapper))
	js.Global().Set("setDevfileContent", js.FuncOf(exports.SetDevfileContentWrapper))
	js.Global().Set("setMetadata", js.FuncOf(exports.SetMetadataWrapper))
	js.Global().Set("updateContainer", js.FuncOf(exports.UpdateContainerWrapper))

	<-make(chan bool)
}
