# gluafs

filesystem utility for [gopher-lua](https://github.com/yuin/gopher-lua). This project is inspired by [layeh/gopher-lfs](https://github.com/layeh/gopher-lfs).

## Installation

```
go get github.com/kohkimakimoto/gluafs
```

## API

### `fs.exists(file)`

### `fs.read(file)`

### `fs.write(file, content)`

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

## Author

Kohki Makimoto <kohki.makimoto@gmail.com>

## License

MIT license.
