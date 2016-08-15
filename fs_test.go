package gluafs

import (
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	"os"
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

var tmpFileContent = "aaaaaaaabbbbbbbb"

func TestRead(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if err = ioutil.WriteFile(tmpFile.Name(), []byte(tmpFileContent), 0644); err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)
	if err := L.DoString(`
local fs = require("fs")
local content = fs.read("` + tmpFile.Name() + `")

assert(content == "aaaaaaaabbbbbbbb")

local content2, err = fs.read("` + tmpFile.Name() + `.hoge")
assert(content2 == nil)


	`); err != nil {
		t.Error(err)
	}
}
