package gluafs

import (
	"io/ioutil"
	"os"

	"fmt"
	"github.com/yookoala/realpath"
	"github.com/yuin/gopher-lua"
	"path/filepath"
)

func Loader(L *lua.LState) int {
	tb := L.NewTable()
	L.SetFuncs(tb, map[string]lua.LGFunction{
		"exists":   exists,
		"read":     read,
		"write":    write,
		"mkdir":    mkdir,
		"remove":   remove,
		"symlink":  symlink,
		"dirname":  dirname,
		"basename": basename,
		"realpath": fnRealpath,
		"getcwd":   getcwd,
		"chdir":    chdir,
		"file":     file,
		"dir":      dir,
		"glob":     glob,
	})
	L.Push(tb)

	return 1
}

func exists(L *lua.LState) int {
	var ret bool

	path := L.CheckString(1)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ret = false
	} else {
		ret = true
	}

	L.Push(lua.LBool(ret))

	return 1
}

func read(L *lua.LState) int {
	path := L.CheckString(1)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		L.Push(lua.LNil)
		return 1
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(string(content)))
	return 1
}

func write(L *lua.LState) int {
	p := L.CheckString(1)
	content := []byte(L.CheckString(2))

	err := ioutil.WriteFile(p, content, os.ModePerm)
	if err != nil {
		L.Push(lua.LBool(false))
		return 1
	}

	L.Push(lua.LBool(true))
	return 1
}

func mkdir(L *lua.LState) int {
	dir := L.CheckString(1)

	perm := 0777
	if L.GetTop() >= 2 {
		perm = L.ToInt(2)
	}

	recursive := false
	if L.GetTop() >= 3 {
		recursive = L.ToBool(3)
	}

	var err error
	if recursive {
		err = os.MkdirAll(dir, os.FileMode(perm))
	} else {
		err = os.Mkdir(dir, os.FileMode(perm))
	}

	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func remove(L *lua.LState) int {
	p := L.CheckString(1)

	recursive := false
	if L.GetTop() >= 2 {
		recursive = L.ToBool(2)
	}

	var err error
	if recursive {
		err = os.Remove(p)
	} else {
		err = os.RemoveAll(p)
	}

	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LTrue)

	return 1
}

func symlink(L *lua.LState) int {
	target := L.CheckString(1)
	link := L.CheckString(2)

	err := os.Symlink(target, link)
	if err != nil {
		panic(err)
	}

	L.Push(lua.LTrue)
	return 1
}

func dirname(L *lua.LState) int {
	filep := L.CheckString(1)
	dirna := filepath.Dir(filep)
	L.Push(lua.LString(dirna))

	return 1
}

func basename(L *lua.LState) int {
	filep := L.CheckString(1)
	dirna := filepath.Base(filep)
	L.Push(lua.LString(dirna))

	return 1
}

func getcwd(L *lua.LState) int {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	L.Push(lua.LString(dir))

	return 1
}

func fnRealpath(L *lua.LState) int {
	filep := L.CheckString(1)

	real, err := realpath.Realpath(filep)
	if err != nil {
		panic(err)
	}
	L.Push(lua.LString(real))

	return 1
}

func chdir(L *lua.LState) int {
	dir := L.CheckString(1)

	err := os.Chdir(dir)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LTrue)
	return 1
}

func file(L *lua.LState) int {
	// same: debug.getinfo(2,'S').source
	var dbg *lua.Debug
	var err error
	var ok bool

	dbg, ok = L.GetStack(1)
	if !ok {
		fmt.Println(dbg)
		L.Push(lua.LNil)
		return 1
	}
	_, err = L.GetInfo("S", dbg, lua.LNil)
	if err != nil {
		fmt.Println(err)
		L.Push(lua.LNil)
		return 1
	}

	L.Push(lua.LString(dbg.Source))
	return 1
}

func dir(L *lua.LState) int {
	// same: debug.getinfo(2,'S').source
	var dbg *lua.Debug
	var err error
	var ok bool

	dbg, ok = L.GetStack(1)
	if !ok {
		fmt.Println(dbg)
		L.Push(lua.LNil)
		return 1
	}
	_, err = L.GetInfo("S", dbg, lua.LNil)
	if err != nil {
		fmt.Println(err)
		L.Push(lua.LNil)
		return 1
	}

	dirname := filepath.Dir(dbg.Source)
	L.Push(lua.LString(dirname))

	return 1
}

func glob(L *lua.LState) int {
	ptn := L.CheckString(1)
	fn := L.CheckFunction(2)

	files, err := filepath.Glob(ptn)
	if err != nil {
		L.RaiseError("Invalid pattern: " + ptn)
		return 0
	}

	for _, f := range files {
		tb := L.NewTable()
		tb.RawSetString("path", lua.LString(f))
		abspath, err := filepath.Abs(f)
		if err != nil {
			L.RaiseError("Invalid path: " + f)
		}
		tb.RawSetString("realpath", lua.LString(abspath))

		err = L.CallByParam(lua.P{
			Fn:      fn,
			NRet:    0,
			Protect: true,
		}, tb)
		if err != nil {
			panic(err)
		}
	}

	return 0
}

func isDir(path string) (ret bool) {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fi.IsDir()
}