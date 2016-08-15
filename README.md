# gluafs

filesystem utility for [gopher-lua](https://github.com/yuin/gopher-lua). This project is inspired by [layeh/gopher-lfs](https://github.com/layeh/gopher-lfs).

## Installation

```
go get github.com/kohkimakimoto/gluafs
```

## API

### `fs.exists(file)`

Returns true if the file exists.

### `fs.read(file)`

Reads the file content and return it. If this function fails, it returns `nil`, plus a string describing the error.

### `fs.write(file, content, [mode])`

Writes the content to the file. If this function fails, it returns `nil`, plus a string describing the error.

### `fs.mkdir(path, recursive)`

### `fs.remove(path, recursive)`

### `fs.symlink(target, link)`

### `fs.dirname(file)`

### `fs.basename(file)`

### `fs.realpath(file)`

### `fs.getcwd()`

### `fs.chdir(path)`

### `fs.file()`

### `fs.dir()`

### `fs.glob(pattern, function)`

## Usage

```go
package main

import (
    "github.com/yuin/gopher-lua"
    "github.com/kohkimakimoto/gluafs"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    L.PreloadModule("fs", gluafs.Loader)
    if err := L.DoString(`
local fs = require("fs")
local ret = fs.exists("path/to/file")

`); err != nil {
        panic(err)
    }
}
```

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

MIT license.
