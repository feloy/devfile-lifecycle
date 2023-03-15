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
	js.Global().Set("setDevfileContent", js.FuncOf(exports.SetDevfileContentWrapper))
	js.Global().Set("setMetadata", js.FuncOf(exports.SetMetadataWrapper))
	js.Global().Set("getFlowChart", js.FuncOf(exports.GetFlowChartWrapper))

	<-make(chan bool)
}
