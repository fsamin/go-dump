package dump_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fsamin/go-dump"
)

func TestDumpStruct(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `T.A: 23
T.B: foo bar
`
	assert.Equal(t, expected, out.String())
}

func TestDumpStructWithPrefix(t *testing.T) {
	type T struct {
		A int
		B string
		C []string
	}

	a := T{23, "foo bar", []string{"c1", "c2"}}

	out := &bytes.Buffer{}
	e := dump.NewEncoder(out)
	e.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultFormatter()}
	e.Prefix = "test"

	err := e.Fdump(a)
	assert.NoError(t, err)

	expected := `test.T.A: 23
test.T.B: foo bar
test.T.C.C0: c1
test.T.C.C1: c2
`
	assert.Equal(t, expected, out.String())
}

func TestDumpArrayWithPrefix(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := []T{{23, "foo bar"}}

	out := &bytes.Buffer{}
	e := dump.NewEncoder(out)
	e.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultFormatter()}
	e.Prefix = "test"

	err := e.Fdump(a)
	assert.NoError(t, err)

	expected := `test.test0.A: 23
test.test0.B: foo bar
`
	assert.Equal(t, expected, out.String())
}

func TestDumpStructWithOutPrefix(t *testing.T) {
	type T struct {
		A  int
		B  string
		TT struct {
			C string
			D string
		}
	}

	a := T{A: 23, B: "foo bar"}
	a.TT.C = "c"
	a.TT.D = "d"

	out := &bytes.Buffer{}
	dumper := dump.NewEncoder(out)
	dumper.DisableTypePrefix = true
	err := dumper.Fdump(a)
	assert.NoError(t, err)

	expected := `A: 23
B: foo bar
TT.C: c
TT.D: d
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
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `T.A: 23
T.B: foo bar
T.C.Cbis: lol
T.C.Cter: lol
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
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `TP.A: 23
TP.B: foo bar
TP.C.Cbis: lol
TP.C.Cter: lol
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
	a.C["bar"] = Tbis{"lel", "lel"}
	a.C["foo"] = Tbis{"lol", "lol"}

	out := &bytes.Buffer{}
	dumper := dump.NewEncoder(out)
	dumper.ExtraFields.Len = true
	dumper.ExtraFields.Type = true

	err := dumper.Fdump(a)
	assert.NoError(t, err)

	expected := `TM.A: 23
TM.B: foo bar
TM.C.__Len__: 2
TM.C.__Type__: Map
TM.C.bar.Tbis.Cbis: lel
TM.C.bar.Tbis.Cter: lel
TM.C.bar.Tbis.__Type__: Tbis
TM.C.foo.Tbis.Cbis: lol
TM.C.foo.Tbis.Cter: lol
TM.C.foo.Tbis.__Type__: Tbis
__Type__: TM
`
	assert.Equal(t, expected, out.String())

}

func TestDumpArray(t *testing.T) {
	a := []T{
		{23, "foo bar", Tbis{"lol", "lol"}},
		{24, "fee bor", Tbis{"lel", "lel"}},
	}

	out := &bytes.Buffer{}
	dumper := dump.NewEncoder(out)
	dumper.ExtraFields.Len = true
	dumper.ExtraFields.Type = true
	err := dumper.Fdump(a)
	assert.NoError(t, err)

	expected := `0.A: 23
0.B: foo bar
0.C.Cbis: lol
0.C.Cter: lol
0.C.__Type__: Tbis
0.__Type__: T
1.A: 24
1.B: fee bor
1.C.Cbis: lel
1.C.Cter: lel
1.C.__Type__: Tbis
1.__Type__: T
__Len__: 2
__Type__: Array
`
	assert.Equal(t, expected, out.String())
}

func TestDumpArray2(t *testing.T) {
	a := []T{
		{23, "foo bar", Tbis{"lol", "lol"}},
		{24, "fee bor", Tbis{"lel", "lel"}},
	}

	out := &bytes.Buffer{}
	dumper := dump.NewEncoder(out)
	dumper.ArrayJSONNotation = true
	dumper.ExtraFields.Len = false
	dumper.ExtraFields.DetailedStruct = false
	dumper.ExtraFields.Type = false
	err := dumper.Fdump(a)
	assert.NoError(t, err)

	expected := `[0].A: 23
[0].B: foo bar
[0].C.Cbis: lol
[0].C.Cter: lol
[1].A: 24
[1].B: fee bor
[1].C.Cbis: lel
[1].C.Cter: lel
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
	dumper := dump.NewEncoder(out)
	dumper.ExtraFields.Len = true
	dumper.ExtraFields.Type = true
	err := dumper.Fdump(a)
	assert.NoError(t, err)
	expected := `TS.A: 0
TS.B: here
TS.C.C0.A: 23
TS.C.C0.B: foo bar
TS.C.C0.C.Cbis: lol
TS.C.C0.C.Cter: lol
TS.C.C0.C.__Type__: Tbis
TS.C.C0.__Type__: T
TS.C.C1.A: 24
TS.C.C1.B: fee bor
TS.C.C1.C.Cbis: lel
TS.C.C1.C.Cter: lel
TS.C.C1.C.__Type__: Tbis
TS.C.C1.__Type__: T
TS.C.__Len__: 2
TS.C.__Type__: Array
TS.D.D0: true
TS.D.D1: false
TS.D.__Len__: 2
TS.D.__Type__: Array
__Type__: TS
`
	assert.Equal(t, expected, out.String())
}

func TestDumpStruct_Array_New_Array_Notation(t *testing.T) {
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
	dumper := dump.NewEncoder(out)
	dumper.ArrayJSONNotation = true
	err := dumper.Fdump(a)
	assert.NoError(t, err)
	expected := `TS.A: 0
TS.B: here
TS.C[0].A: 23
TS.C[0].B: foo bar
TS.C[0].C.Cbis: lol
TS.C[0].C.Cter: lol
TS.C[1].A: 24
TS.C[1].B: fee bor
TS.C[1].C.Cbis: lel
TS.C[1].C.Cter: lel
TS.D[0]: true
TS.D[1]: false
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
		t.Logf("%s: %v (%T)", k, v, v)
		if k == "T.A" {
			m1Found = true
			assert.Equal(t, 23, v)
		}
		if k == "T.B" {
			m2Found = true
			assert.Equal(t, "foo bar", v)
		}
	}
	assert.True(t, m1Found, "T.A not found in map")
	assert.True(t, m2Found, "T.B not found in map")
}

func TestToMapWithFormatter(t *testing.T) {
	type T struct {
		A int
		B string
	}

	a := T{23, "foo bar"}

	m, err := dump.ToMap(a, dump.WithDefaultLowerCaseFormatter())
	t.Log(m)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(m))
	var m1Found, m2Found bool
	for k, v := range m {
		if k == "t.a" {
			m1Found = true
			assert.Equal(t, 23, v)
		}
		if k == "t.b" {
			m2Found = true
			assert.Equal(t, "foo bar", v)
		}
	}
	assert.True(t, m1Found, "t.a not found in map")
	assert.True(t, m2Found, "t.b not found in map")
}

func TestMapStringInterface(t *testing.T) {
	myMap := make(map[string]interface{})
	myMap["id"] = "ID"
	myMap["name"] = "foo"
	myMap["value"] = "bar"
	myMap[""] = "empty"

	result, err := dump.ToStringMap(myMap)
	t.Log(dump.Sdump(myMap))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(result))

	expected := `id: ID
name: foo
value: bar
`
	out := &bytes.Buffer{}
	err = dump.Fdump(out, myMap, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}

func TestMapEmptyInterface(t *testing.T) {
	myMap := make(map[string]interface{})
	myMap[""] = "empty"

	result, err := dump.ToStringMap(myMap)
	t.Log(dump.Sdump(myMap))
	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))

	expected := ``
	out := &bytes.Buffer{}
	err = dump.Fdump(out, myMap, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}

func TestFromJSON(t *testing.T) {
	js := []byte(`{
    "blabla": "lol log", 
    "boubou": {
        "yo": 1
    } 
}`)

	var i interface{}
	assert.NoError(t, json.Unmarshal(js, &i))

	result, err := dump.ToStringMap(i)
	t.Log(dump.Sdump(i))
	t.Log(result)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "lol log", result["blabla"])
	assert.Equal(t, "1", result["boubou.yo"])
}

type Result struct {
	Body     string      `json:"body,omitempty" yaml:"body,omitempty"`
	BodyJSON interface{} `json:"bodyjson,omitempty" yaml:"bodyjson,omitempty"`
}

func TestMapStringInterfaceInStruct(t *testing.T) {
	r := Result{}
	r.Body = "foo"
	r.BodyJSON = map[string]interface{}{
		"cardID": "1234",
		"items":  []string{"foo", "beez"},
		"test": Result{
			Body: "12",
			BodyJSON: map[string]interface{}{
				"card": "@",
				"yolo": 3,
				"beez": true,
			},
		},
		"description": "yolo",
	}

	expected := `result.body: foo
result.bodyjson.cardid: 1234
result.bodyjson.description: yolo
result.bodyjson.items.items0: foo
result.bodyjson.items.items1: beez
result.bodyjson.test.result.body: 12
result.bodyjson.test.result.bodyjson.beez: true
result.bodyjson.test.result.bodyjson.card: @
result.bodyjson.test.result.bodyjson.yolo: 3
`

	out := &bytes.Buffer{}
	err := dump.Fdump(out, r, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}

func TestWeird(t *testing.T) {
	testJSON := `{
	"beez": null,
	"foo" : "bar",
	"bou" : [null, "hello"]
  }`

	var test interface{}
	json.Unmarshal([]byte(testJSON), &test)
	expected := `beez:
bou.bou0:
bou.bou1: hello
foo: bar
`

	out := &bytes.Buffer{}
	err := dump.Fdump(out, test, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())

}

type ResultUnexported struct {
	body *string
	Foo  string
}

func TestUnexportedField(t *testing.T) {

	test := ResultUnexported{
		body: nil,
		Foo:  "bar",
	}

	expected := `resultunexported.foo: bar
`

	out := &bytes.Buffer{}
	err := dump.Fdump(out, test, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())
}

func TestWithDetailedStruct(t *testing.T) {
	type T struct {
		A int
		B string
		T *T
	}

	a := T{23, "foo bar", &T{A: 46, B: "fiz buzz"}}

	enc := dump.NewDefaultEncoder()
	enc.ExtraFields.DetailedStruct = true
	enc.ExtraFields.Type = false
	enc.ExtraFields.Len = true
	res, _ := enc.Sdump(a)
	t.Log(res)
	assert.Equal(t, `T.A: 23
T.B: foo bar
T.T: {"A":46,"B":"fiz buzz","T":null}
T.T.A: 46
T.T.B: fiz buzz
T.T.T: 
T.T.__Len__: 3
T.__Len__: 3
`, res)
}

func TestDumpJSONInString(t *testing.T) {
	type T struct {
		A int
		B string
	}
	value := T{
		A: 0,
		B: "{ \"toctoc\": \"Qui est la\"}",
	}

	e := dump.NewDefaultEncoder()
	e.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultLowerCaseFormatter()}
	e.ExtraFields.DetailedMap = false
	e.ExtraFields.DetailedStruct = false
	e.ExtraFields.DeepJSON = true
	e.ExtraFields.Len = false
	e.ExtraFields.Type = false
	m, err := e.ToStringMap(value)
	assert.NoError(t, err)
	assert.Equal(t, "Qui est la", m["t.b.toctoc"])
}

func TestNoDumpJSONInString(t *testing.T) {
	type T struct {
		A int
		B string
	}
	value := T{
		A: 0,
		B: "{ \"toctoc\": \"Qui est la\"}",
	}

	e := dump.NewDefaultEncoder()
	e.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultLowerCaseFormatter()}
	e.ExtraFields.DetailedMap = false
	e.ExtraFields.DetailedStruct = false
	e.ExtraFields.DeepJSON = false
	e.ExtraFields.Len = false
	e.ExtraFields.Type = false
	m, err := e.ToStringMap(value)
	assert.NoError(t, err)
	assert.Equal(t, "{ \"toctoc\": \"Qui est la\"}", m["t.b"])
}

func TestBuildEnvironmentVariable(t *testing.T) {
	type B struct {
		PartOfB string
	}
	type PieceOfConfig struct {
		B
		C string
	}
	type Config struct {
		A    string
		MapB map[string]PieceOfConfig
	}

	cfg := Config{
		A: "A",
		MapB: map[string]PieceOfConfig{
			"1": {
				B: B{PartOfB: "B"},
				C: "C",
			},
			"2": {
				B: B{PartOfB: "2B"},
				C: "2C",
			},
		},
	}

	out := &bytes.Buffer{}
	dumper := dump.NewEncoder(out)
	dumper.DisableTypePrefix = true
	dumper.Separator = "_"
	dumper.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultUpperCaseFormatter()}
	err := dumper.Fdump(cfg)
	assert.NoError(t, err)
	expected := `A: A
MAPB_1_B_PARTOFB: B
MAPB_1_C: C
MAPB_2_B_PARTOFB: 2B
MAPB_2_C: 2C
`
	assert.Equal(t, expected, out.String())

	out = &bytes.Buffer{}
	dumper = dump.NewEncoder(out)
	dumper.DisableTypePrefix = false
	dumper.Separator = "_"
	dumper.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultUpperCaseFormatter()}
	err = dumper.Fdump(cfg)
	assert.NoError(t, err)
	expected = `CONFIG_A: A
CONFIG_MAPB_1_PIECEOFCONFIG_B_PARTOFB: B
CONFIG_MAPB_1_PIECEOFCONFIG_C: C
CONFIG_MAPB_2_PIECEOFCONFIG_B_PARTOFB: 2B
CONFIG_MAPB_2_PIECEOFCONFIG_C: 2C
`
	assert.Equal(t, expected, out.String())
}

func TestEnvVariableWithViper(t *testing.T) {
	type MyStruct struct {
		A string
		B struct {
			InsideB string
		}
		/*C struct {
			InnerC
			InsideC string
		}
		D map[string]D*/
	}

	var myStruct MyStruct
	myStruct.A = "value A"
	myStruct.B.InsideB = "value B"
	//myStruct.C.InsideInnerC = "value C (inner C)"
	//myStruct.C.InsideC = "value C"
	//myStruct.D = map[string]D{
	//	"d1": D{InsideD: "value D1"},
	//	"d2": D{InsideD: "value D2"},
	//}

	dumper := dump.NewDefaultEncoder()
	dumper.DisableTypePrefix = true
	dumper.Separator = "_"
	dumper.Prefix = "MYSTRUCT"
	dumper.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultUpperCaseFormatter()}
	envs, err := dumper.ToStringMap(&myStruct)
	assert.NoError(t, err)
	t.Log("go-dump result...")

	for k, v := range envs {
		t.Log(k, v)
		os.Setenv(k, v)
		viper.BindEnv(dumper.ViperKey(k), k)
	}

	t.Log("env...")
	for _, e := range os.Environ() {
		t.Log(e)
	}

	viperSettings := viper.AllSettings()
	t.Log("viper results...")

	for k, v := range viperSettings {
		t.Log(k, v)
	}

	t.Log("viper unmarshal...")

	var myStructFromViper MyStruct
	err = viper.Unmarshal(&myStructFromViper)
	assert.NoError(t, err)
	t.Logf("--> %T %+v", myStructFromViper, myStructFromViper)

	assert.Equal(t, myStruct, myStructFromViper)

}

func TestTruc(t *testing.T) {
	var b = `{
		"user": "STRING",
		"password": "STRING",
		"database": "STRING",
		"writers" : [
			{ "host": "STRING", "port": "-1" }
		],
		"readers" : [
			{ "host": "STRING", "port": "-1"  }
		],
		"analytics" : [
			{ "host": "STRING", "port": "-1"  }
		],
		"type": "STRING",
		"ssl": "(off|preferred|required|strict)"
	 }`

	var x = map[string]interface{}{}
	assert.NoError(t, json.Unmarshal([]byte(b), &x))

	iDumped, err := dump.ToStringMap(x)
	assert.NoError(t, err)

	fmt.Println(iDumped)
}

func TestDetailArray(t *testing.T) {

	var b = `{
		"user": "STRING",
		"password": "STRING",
		"database": "STRING",
		"writers" : [
			{ "host": "STRING", "port": "-1" }
		],
		"readers" : [
			{ "host": "STRING", "port": "-1"  }
		],
		"analytics" : [
			{ "host": "STRING", "port": "-1"  }
		],
		"type": "STRING",
		"ssl": "(off|preferred|required|strict)"
	 }`

	var x = map[string]interface{}{}
	assert.NoError(t, json.Unmarshal([]byte(b), &x))

	e := dump.NewDefaultEncoder()
	e.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultLowerCaseFormatter()}
	e.ExtraFields.DetailedMap = true
	e.ExtraFields.DetailedStruct = true
	e.ExtraFields.DetailedArray = true

	m, err := e.ToMap(x)
	assert.NoError(t, err)
	t.Logf("shoud be an array: %T %v", m["writers"], m["writers"])
	assert.NotEmpty(t, m["writers"])
}

func Test_DumpTime(t *testing.T) {
	m := map[string]interface{}{
		"string": "foobar",
		"date":   time.Date(2020, time.November, 29, 10, 00, 00, 00, time.Local),
		"dates":  []time.Time{time.Date(2020, time.November, 29, 10, 00, 00, 00, time.Local)},
	}

	result, err := dump.ToStringMap(m)
	require.NoError(t, err)
	t.Log(result)

	require.NotZero(t, result["date"])
	require.NotZero(t, result["date.Time"])
	require.NotZero(t, result["dates.dates0"])

}

func Test_DumpMap(t *testing.T) {
	m := map[string]interface{}{
		"string": "foobar",
	}
	e := dump.NewDefaultEncoder()
	e.ExtraFields.Len = true
	e.ExtraFields.Type = true
	e.ExtraFields.DetailedStruct = true
	e.ExtraFields.DetailedMap = true
	e.ExtraFields.DetailedArray = true
	result, err := e.ToStringMap(m)
	require.NoError(t, err)
	t.Log(result)
	require.Len(t, result, 3)
}

func Test_DumpTimeWithDetailledStruct(t *testing.T) {
	m := map[string]interface{}{
		"string": "foobar",
		"date":   time.Date(2020, time.November, 29, 10, 00, 00, 00, time.Local),
		"dates":  []time.Time{time.Date(2020, time.November, 29, 10, 00, 00, 00, time.Local)},
	}

	dmp := dump.NewDefaultEncoder()
	dmp.ExtraFields.DetailedStruct = true
	dmp.ExtraFields.DetailedMap = true
	dmp.ExtraFields.DetailedArray = true
	result, err := dmp.ToStringMap(m)
	t.Log(result)
	require.NoError(t, err)
	require.NotZero(t, result["date"])
	require.NotZero(t, result["date.Time"])
	require.NotZero(t, result["dates.dates0"])

}

func Test_DumpResultStruct(t *testing.T) {
	type Result struct {
		Foo string
		Bar string
	}
	m := Result{Foo: "foo", Bar: "bar"}
	e := dump.NewDefaultEncoder()
	e.ExtraFields.Len = true
	e.ExtraFields.Type = true
	e.ExtraFields.DetailedStruct = true
	e.ExtraFields.DetailedMap = true
	e.ExtraFields.DetailedArray = true
	result, err := e.ToStringMap(m)
	require.NoError(t, err)
	t.Log(result)
	require.Len(t, result, 4)
}

func Test_DumpArrayResultStruct(t *testing.T) {
	type Result struct {
		Foo string
		Bar string
	}
	m := []Result{
		{Foo: "foo1", Bar: "bar1"},
		{Foo: "foo2", Bar: "bar2"},
	}
	e := dump.NewDefaultEncoder()
	e.ExtraFields.Len = true
	e.ExtraFields.Type = true
	e.ExtraFields.DetailedStruct = true
	e.ExtraFields.DetailedMap = true
	e.ExtraFields.DetailedArray = true
	result, err := e.ToStringMap(m)
	require.NoError(t, err)
	t.Log(result)
	require.Len(t, result, 10)
}

func Test_DumpJSONAnnotationResultStruct(t *testing.T) {
	type Result struct {
		Foo  string `json:"Foo2,omitempty"`
		Bar  string `json:",omitempty"`
		Bar2 string `json:",omitempty"`
	}
	m := Result{Foo: "foo1", Bar: "bar", Bar2: "bar2"}

	e := dump.NewDefaultEncoder()
	e.ExtraFields.Len = false
	e.ExtraFields.Type = false
	e.ExtraFields.DetailedStruct = true
	e.ExtraFields.DetailedMap = true
	e.ExtraFields.DetailedArray = true

	e.ExtraFields.UseJSONTag = true

	e.Formatters = []dump.KeyFormatterFunc{
		func(s string, level int) string {
			if level == 0 {
				return strings.ToLower(s)
			}
			return s
		},
	}

	result, err := e.ToStringMap(m)

	require.NoError(t, err)
	t.Log(result)
	require.Equal(t, "foo1", result["result.Foo2"])
	require.Equal(t, "bar", result["result.Bar"])
	require.Equal(t, "bar2", result["result.Bar2"])
}

func TestDumpStructWithPrefixJson(t *testing.T) {
	type PieceOfConfig struct {
		C string `json:"ZC`
	}

	type T struct {
		A  int `json:"AX"`
		B  string
		TT struct {
			C string
			D string `json:"DY"`
		}
		MapB map[string]PieceOfConfig `json:"MMapB"`
	}

	a := T{A: 23, B: "foo bar"}
	a.TT.C = "c"
	a.TT.D = "d"

	a.MapB = map[string]PieceOfConfig{
		"1": {
			C: "C",
		},
	}

	out := &bytes.Buffer{}
	dumper := dump.NewEncoder(out)
	dumper.DisableTypePrefix = true
	dumper.ExtraFields.UseJSONTag = true
	err := dumper.Fdump(a)
	assert.NoError(t, err)

	expected := `AX: 23
B: foo bar
MMapB.1.C: C
TT.C: c
TT.DY: d
`
	assert.Equal(t, expected, out.String())
}
