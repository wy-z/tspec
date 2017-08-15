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
package tests

import "time"

// A A ...
type A struct {
    Bool       bool
    Uint       uint
    Uint8      uint8
    Uint16     uint16
    Uint32     uint32
    Uint64     uint64
    Int        int
    Int8       int8
    Int16      int16
    Int32      int32
    Int64      int64
    Uintptr    uintptr
    Float32    float32
    Float64    float64
    String     string
    Complex64  complex64
    Complex128 complex128
    Byte       byte
    Rune       rune
    Time       time.Time
}

// B B ...
type B struct {
    BStruct *B
}

// C C ...
type C struct {
    CStruct *struct {
        String string
        Bool   bool
    }
}

// D D ...
type D struct {
    DArray []string
}

// E E ...
type E struct {
    EArray []*C
}

// F F ...
type F struct {
    FMap map[string]string
}

// G G ...
type G struct {
    *C
}
```

run ```tspec A```

```
{
    "tests.A": {
        "id": "tests.A",
        "type": "object",
        "title": "A",
        "properties": {
            "Bool": {
                "type": "boolean"
            },
            "Byte": {
                "type": "string",
                "format": "byte"
            },
            "Complex128": {
                "type": "number",
                "format": "double"
            },
            "Complex64": {
                "type": "number",
                "format": "float"
            },
            "Float32": {
                "type": "number",
                "format": "float32"
            },
            "Float64": {
                "type": "number",
                "format": "float64"
            },
            "Int": {
                "type": "integer",
                "format": "int"
            },
            "Int16": {
                "type": "integer",
                "format": "int16"
            },
            "Int32": {
                "type": "integer",
                "format": "int32"
            },
            "Int64": {
                "type": "integer",
                "format": "int64"
            },
            "Int8": {
                "type": "integer",
                "format": "int8"
            },
            "Rune": {
                "type": "string",
                "format": "byte"
            },
            "String": {
                "type": "string"
            },
            "Time": {
                "type": "string",
                "format": "date-time"
            },
            "Uint": {
                "type": "integer",
                "format": "uint"
            },
            "Uint16": {
                "type": "integer",
                "format": "uint16"
            },
            "Uint32": {
                "type": "integer",
                "format": "uint32"
            },
            "Uint64": {
                "type": "integer",
                "format": "uint64"
            },
            "Uint8": {
                "type": "integer",
                "format": "uint8"
            },
            "Uintptr": {
                "type": "integer",
                "format": "int64"
            }
        }
    }
}
```
