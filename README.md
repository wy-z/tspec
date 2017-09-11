# tspec
[![Build Status](https://travis-ci.org/wy-z/tspec.svg?branch=master)](https://travis-ci.org/wy-z/tspec) [![GoDoc](https://godoc.org/github.com/wy-z/tspec?status.svg)](http://godoc.org/github.com/wy-z/tspec) [![Go Report Card](https://goreportcard.com/badge/github.com/wy-z/tspec)](https://goreportcard.com/report/github.com/wy-z/tspec)

Parse golang data structure into json schema.

## Installation
```
go get github.com/wy-z/tspec
```
Or
```
import "github.com/wy-z/tspec/tspec" # see cli/cli.go
```

## Usage
```
NAME:
   TSpec - Parse golang data structure into json schema.

USAGE:
   tspec [global options] command [command options] [arguments...]

VERSION:
   1.9.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --package PKG, -p PKG           package dir PKG (default: ".")
   --expression EXPR, --expr EXPR  type expression EXPR
   --help, -h                      show help
   --version, -v                   print the version
```

## QuickStart

```
package main

import (
    "github.com/wy-z/tspec/samples"
)

func main() {
    _ = new(samples.NormalStruct)
}
```

run ```tspec samples.NormalStruct```

## Samples

see `github.com/wy-z/tspec/samples/source`

## Test

```
go get -u github.com/jteeuwen/go-bindata/...
go generate ./samples && go test -v ./tspec
```