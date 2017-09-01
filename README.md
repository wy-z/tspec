# tspec
[![Build Status](https://travis-ci.org/wy-z/tspec.svg?branch=master)](https://travis-ci.org/wy-z/tspec) [![GoDoc](https://godoc.org/github.com/wy-z/tspec?status.svg)](http://godoc.org/github.com/wy-z/tspec)

Parse golang data structure into json schema.

## Installation
```
go get github.com/wy-z/tspec
```
Or
```
import "github.com/wy-z/tspec/tspec" # see main.go
```

## Usage
```
NAME:
   TSpec - Parse golang data structure into json schema.

USAGE:
   tspec [global options] command [command options] [arguments...]

VERSION:
   1.2.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --package PKG, -p PKG           package dir PKG (default: ".")
   --expression EXPR, --expr EXPR  type expression EXPR
   --help, -h                      show help
   --version, -v                   print the version
```

## QuickStart

type_samples.go

```
package tspec_test

import (
    "github.com/wy-z/tspec/samples"
)
```

run ```tspec NormalStruct```

## Samples

see `github.com/wy-z/tspec/samples`

## Test

```
go get github.com/jteeuwen/go-bindata
go-bindata -o samples/samples.go  -pkg samples samples/source
go test -v ./tspec
```