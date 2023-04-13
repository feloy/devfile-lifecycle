package exports

import (
	"fmt"
	"syscall/js"

	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	"github.com/feloy/devfile-lifecycle/client/devfile-builder/wasm/pkg/global"
)

func MoveCommandWrapper(this js.Value, args []js.Value) interface{} {
	return result(
		moveCommand(args[0].String(), args[1].String(), args[2].Int(), args[3].Int()),
	)
}

func moveCommand(previousKind, newKind string, previousIndex, newIndex int) (map[string]interface{}, error) {
	commands, err := global.Devfile.Data.GetCommands(common.DevfileOptions{})
	if err != nil {
		return nil, err
	}

	currentPreviousIndex := 0
	var previousGlobalIndex int
	for i, command := range commands {
		if getGroup(command) == previousKind {
			if currentPreviousIndex == previousIndex {
				previousGlobalIndex = i
				fmt.Printf("found command %s\n", command.Id)
				setGroup(&command, newKind)
				err = global.Devfile.Data.UpdateCommand(command)
				if err != nil {
					return nil, err
				}
				break
			}
			currentPreviousIndex++
		}
	}
	fmt.Printf("index in global list: %d\n", previousGlobalIndex)
	return getContent()
}
