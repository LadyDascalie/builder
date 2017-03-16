# builder : A very naive multi-platform `go build` 

## Overview [![GoDoc](https://godoc.org/github.com/ladydascalie/builder?status.svg)](https://godoc.org/github.com/ladydascalie/builder) [![Go Report Card](https://goreportcard.com/badge/github.com/ladydascalie/builder)](https://goreportcard.com/report/github.com/ladydascalie/builder)

Builder will run go build for macOS, Linux and Windows, for adm64 and 386 architectures.
Once it is done running it will restore your `GOARCH` and `GOOS` environment variables to their values before `builder` ran.

Usages: 
```
# Just run it naked to cross-compile for macOS, Linux and Windows:
builder

# Or pick a system to build for:
builder -for linux

```

## Install

```
go get github.com/ladydascalie/builder
```

## License

MIT.
