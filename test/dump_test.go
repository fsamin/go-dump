package test

import (
	"testing"

	"bytes"

	"github.com/fsamin/go-dump"
	"github.com/stretchr/testify/assert"
)

func TestDumpStruct(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	out := &bytes.Buffer{}
	err := dump.FDump(out, a)
	assert.NoError(t, err)

	expected := `T.A: 23
T.B: foo bar
`
	assert.Equal(t, expected, out.String())

}

type T struct {
	A int
	B string
	C Tbis
}

type Tbis struct {
	Cbis string
	Cter string
}

func TestDumpStruct_Nested(t *testing.T) {

	a := T{23, "foo bar", Tbis{"lol", "lol"}}

	out := &bytes.Buffer{}
	err := dump.FDump(out, a)
	assert.NoError(t, err)

	expected := `T.A: 23
T.B: foo bar
T.C.Tbis.Cbis: lol
T.C.Tbis.Cter: lol
`
	assert.Equal(t, expected, out.String())

}

type TP struct {
	A *int
	B string
	C *Tbis
}

func TestDumpStruct_NestedWithPointer(t *testing.T) {
	i := 23
	a := TP{&i, "foo bar", &Tbis{"lol", "lol"}}

	out := &bytes.Buffer{}
	err := dump.FDump(out, a)
	assert.NoError(t, err)

	expected := `TP.A: 23
TP.B: foo bar
TP.C.Tbis.Cbis: lol
TP.C.Tbis.Cter: lol
`
	assert.Equal(t, expected, out.String())

}

type TM struct {
	A int
	B string
	C map[string]Tbis
}

func TestDumpStruct_Map(t *testing.T) {

	a := TM{A: 23, B: "foo bar"}
	a.C = map[string]Tbis{}
	a.C["1"] = Tbis{"lol", "lol"}
	a.C["2"] = Tbis{"lel", "lel"}

	out := &bytes.Buffer{}
	err := dump.FDump(out, a)
	assert.NoError(t, err)

	expected := `TM.A: 23
TM.B: foo bar
TM.C.1.Tbis.Cbis: lol
TM.C.1.Tbis.Cter: lol
TM.C.2.Tbis.Cbis: lel
TM.C.2.Tbis.Cter: lel
`
	assert.Equal(t, expected, out.String())

}

func TestDumpArray(t *testing.T) {
	a := []T{
		{23, "foo bar", Tbis{"lol", "lol"}},
		{24, "fee bor", Tbis{"lel", "lel"}},
	}

	out := &bytes.Buffer{}
	err := dump.FDump(out, a)
	assert.NoError(t, err)

	expected := `0.T.A: 23
0.T.B: foo bar
0.T.C.Tbis.Cbis: lol
0.T.C.Tbis.Cter: lol
1.T.A: 24
1.T.B: fee bor
1.T.C.Tbis.Cbis: lel
1.T.C.Tbis.Cter: lel
`
	assert.Equal(t, expected, out.String())
}

type TS struct {
	A int
	B string
	C []T
	D []bool
}

func TestDumpStruct_Array(t *testing.T) {
	a := TS{
		A: 0,
		B: "here",
		C: []T{
			{23, "foo bar", Tbis{"lol", "lol"}},
			{24, "fee bor", Tbis{"lel", "lel"}},
		},
		D: []bool{true, false},
	}

	out := &bytes.Buffer{}
	err := dump.FDump(out, a)
	assert.NoError(t, err)
	expected := `TS.A: 0
TS.B: here
TS.C.C0.T.A: 23
TS.C.C0.T.B: foo bar
TS.C.C0.T.C.Tbis.Cbis: lol
TS.C.C0.T.C.Tbis.Cter: lol
TS.C.C1.T.A: 24
TS.C.C1.T.B: fee bor
TS.C.C1.T.C.Tbis.Cbis: lel
TS.C.C1.T.C.Tbis.Cter: lel
TS.D.D0: true
TS.D.D1: false
`
	assert.Equal(t, expected, out.String())
}

func TestToMap(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	m, err := dump.ToMap(a)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(m))
	var m1Found, m2Found bool
	for k, v := range m {
		if k == "T.A" {
			m1Found = true
			assert.Equal(t, "23", v)
		}
		if k == "T.B" {
			m2Found = true
			assert.Equal(t, "foo bar", v)
		}
	}
	assert.True(t, m1Found, "T.A not found in map")
	assert.True(t, m2Found, "T.B not found in map")
}
