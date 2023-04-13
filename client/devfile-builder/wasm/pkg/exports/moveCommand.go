package exports

import (
	"fmt"
	"syscall/js"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/v2/pkg/devfile/parser/data/v2/common"
	"github.com/feloy/devfile-lifecycle/client/devfile-builder/wasm/pkg/global"
)

func MoveCommandWrapper(this js.Value, args []js.Value) interface{} {
	return result(
		moveCommand(args[0].String(), args[1].String(), args[2].Int(), args[3].Int()),
	)
}

func moveCommand(previousGroup, newGroup string, previousIndex, newIndex int) (map[string]interface{}, error) {
	commands, err := global.Devfile.Data.GetCommands(common.DevfileOptions{})
	if err != nil {
		return nil, err
	}

	commandsByGroup := map[string][]v1alpha2.Command{}

	for _, command := range commands {
		group := getGroup(command)
		commandsByGroup[group] = append(commandsByGroup[group], command)
		global.Devfile.Data.DeleteCommand(command.Id)
	}

	if len(commandsByGroup[previousGroup]) < previousIndex {
		return nil, fmt.Errorf("unable to find command at index #%d in group %q", previousIndex, previousGroup)
	}

	commandToMove := commandsByGroup[previousGroup][previousIndex]
	setGroup(&commandToMove, newGroup)

	commandsByGroup[previousGroup] = append(
		commandsByGroup[previousGroup][:previousIndex],
		commandsByGroup[previousGroup][previousIndex+1:]...,
	)

	end := append([]v1alpha2.Command{}, commandsByGroup[newGroup][newIndex:]...)
	commandsByGroup[newGroup] = append(commandsByGroup[newGroup][:newIndex], commandToMove)
	commandsByGroup[newGroup] = append(commandsByGroup[newGroup], end...)

	for _, group := range []string{"build", "run", "test", "debug", "deploy", ""} {
		err = global.Devfile.Data.AddCommands(commandsByGroup[group])
		if err != nil {
			return nil, err
		}
	}

	return getContent()
}
