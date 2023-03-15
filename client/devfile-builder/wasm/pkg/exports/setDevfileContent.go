package exports

import (
	"errors"
	"syscall/js"

	"github.com/devfile/library/v2/pkg/devfile"
	"github.com/devfile/library/v2/pkg/devfile/parser"
	context "github.com/devfile/library/v2/pkg/devfile/parser/context"
	"github.com/devfile/library/v2/pkg/testingutil/filesystem"

	"github.com/feloy/devfile-lifecycle/client/devfile-builder/wasm/pkg/global"
)

// setDevfileContent

func SetDevfileContentWrapper(this js.Value, args []js.Value) interface{} {
	return result(
		SetDevfileContent(args[0].String()),
	)
}

func SetDevfileContent(content string) (map[string]interface{}, error) {
	parserArgs := parser.ParserArgs{
		Data: []byte(content),
	}
	var err error
	global.Devfile, _, err = devfile.ParseDevfileAndValidate(parserArgs)
	if err != nil {
		return nil, errors.New("error parsing devfile")
	}
	global.FS = filesystem.NewFakeFs()
	global.Devfile.Ctx = context.FakeContext(global.FS, "/devfile.yaml")

	return getContent()
}
