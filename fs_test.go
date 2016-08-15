package gluafs

import (
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestExists(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)
	if err := L.DoString(`
local fs = require("fs")
assert(fs.exists(".") == true)

	`); err != nil {
		t.Error(err)
	}
}
