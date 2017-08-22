# tspec
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
   0.0.1

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