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
