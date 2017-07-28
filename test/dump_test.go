package test

import (
	"bytes"
	"testing"

	"github.com/ovh/cds/engine/api/test"
	"github.com/ovh/cds/sdk"
	"github.com/stretchr/testify/assert"

	"encoding/json"

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
	a.C["1"] = Tbis{"lol", "lol"}
	a.C["2"] = Tbis{"lel", "lel"}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `TM.A: 23
TM.B: foo bar
TM.C.1.Tbis.Cbis: lol
TM.C.1.Tbis.Cter: lol
TM.C.2.Tbis.Cbis: lel
TM.C.2.Tbis.Cter: lel
TM.C.Len: 2
`
	assert.Equal(t, expected, out.String())

}

func TestDumpArray(t *testing.T) {
	a := []T{
		{23, "foo bar", Tbis{"lol", "lol"}},
		{24, "fee bor", Tbis{"lel", "lel"}},
	}

	out := &bytes.Buffer{}
	err := dump.Fdump(out, a)
	assert.NoError(t, err)

	expected := `0.A: 23
0.B: foo bar
0.C.Cbis: lol
0.C.Cter: lol
1.A: 24
1.B: fee bor
1.C.Cbis: lel
1.C.Cter: lel
Len: 2
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
	err := dump.Fdump(out, a)
	assert.NoError(t, err)
	expected := `TS.A: 0
TS.B: here
TS.C.C0.A: 23
TS.C.C0.B: foo bar
TS.C.C0.C.Cbis: lol
TS.C.C0.C.Cter: lol
TS.C.C1.A: 24
TS.C.C1.B: fee bor
TS.C.C1.C.Cbis: lel
TS.C.C1.C.Cter: lel
TS.C.Len: 2
TS.D.D0: true
TS.D.D1: false
TS.D.Len: 2
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
			assert.Equal(t, "23", v)
		}
		if k == "t.b" {
			m2Found = true
			assert.Equal(t, "foo bar", v)
		}
	}
	assert.True(t, m1Found, "t.a not found in map")
	assert.True(t, m2Found, "t.b not found in map")
}

func TestComplex(t *testing.T) {
	p := sdk.Pipeline{
		Name: "MyPipeline",
		Type: sdk.BuildPipeline,
		Stages: []sdk.Stage{
			{
				BuildOrder: 1,
				Name:       "stage 1",
				Enabled:    true,
				Jobs: []sdk.Job{
					{
						Action: sdk.Action{
							Name:        "Job 1",
							Description: "This is job 1",
							Actions: []sdk.Action{
								{

									Type: sdk.BuiltinAction,
									Name: sdk.ScriptAction,
									Parameters: []sdk.Parameter{
										{
											Name:  "script",
											Type:  sdk.TextParameter,
											Value: "echo lol",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Log(dump.MustSdump(p))

	out := &bytes.Buffer{}
	err := dump.Fdump(out, p)
	assert.NoError(t, err)
	expected := `Pipeline.AttachedApplication.Len: 0
Pipeline.GroupPermission.Len: 0
Pipeline.ID: 0
Pipeline.LastModified: 0
Pipeline.Name: MyPipeline
Pipeline.Parameter.Len: 0
Pipeline.Permission: 0
Pipeline.ProjectID: 0
Pipeline.Stages.Len: 1
Pipeline.Stages.Stages0.BuildOrder: 1
Pipeline.Stages.Stages0.Enabled: true
Pipeline.Stages.Stages0.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Actions.Len: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Enabled: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Final: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.LastModified: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Name: Script
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Len: 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.Name: script
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.Type: text
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Parameters.Parameters0.Value: echo lol
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Requirements.Len: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Actions0.Type: Builtin
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Actions.Len: 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Description: This is job 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Enabled: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Final: false
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.ID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.LastModified: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Name: Job 1
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Parameters.Len: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Action.Requirements.Len: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.Enabled: false
Pipeline.Stages.Stages0.Jobs.Jobs0.LastModified: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.PipelineActionID: 0
Pipeline.Stages.Stages0.Jobs.Jobs0.PipelineStageID: 0
Pipeline.Stages.Stages0.Jobs.Len: 1
Pipeline.Stages.Stages0.LastModified: 0
Pipeline.Stages.Stages0.Name: stage 1
Pipeline.Stages.Stages0.PipelineBuildJobs.Len: 0
Pipeline.Stages.Stages0.PipelineID: 0
Pipeline.Stages.Stages0.Prerequisites.Len: 0
Pipeline.Stages.Stages0.RunJobs.Len: 0
Pipeline.Type: build
`
	assert.Equal(t, expected, out.String())
	assert.NoError(t, err)
}

func TestMapStringInterface(t *testing.T) {
	myMap := make(map[string]interface{})
	myMap["id"] = "ID"
	myMap["name"] = "foo"
	myMap["value"] = "bar"

	result, err := dump.ToMap(myMap)
	t.Log(dump.Sdump(myMap))
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result))
}

func TestFromJSON(t *testing.T) {
	js := []byte(`{
    "blabla": "lol log",
    "boubou": {
        "yo": 1
    }
}`)

	var i interface{}
	test.NoError(t, json.Unmarshal(js, &i))

	result, err := dump.ToMap(i)
	t.Log(dump.Sdump(i))
	t.Log(result)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(result))
	assert.Equal(t, "lol log", result["blabla"])
	assert.Equal(t, "1", result["boubou.yo"])
}
func TestComplexWithFormatter(t *testing.T) {
	p := sdk.Pipeline{
		Name: "MyPipeline",
		Type: sdk.BuildPipeline,
		Stages: []sdk.Stage{
			{
				BuildOrder: 1,
				Name:       "stage 1",
				Enabled:    true,
				Jobs: []sdk.Job{
					{
						Action: sdk.Action{
							Name:        "Job 1",
							Description: "This is job 1",
							Actions: []sdk.Action{
								{

									Type: sdk.BuiltinAction,
									Name: sdk.ScriptAction,
									Parameters: []sdk.Parameter{
										{
											Name:  "script",
											Type:  sdk.TextParameter,
											Value: "echo lol",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	t.Log(dump.MustSdump(p))

	out := &bytes.Buffer{}
	err := dump.Fdump(out, p, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	expected := `pipeline.attachedapplication.len: 0
pipeline.grouppermission.len: 0
pipeline.id: 0
pipeline.lastmodified: 0
pipeline.name: MyPipeline
pipeline.parameter.len: 0
pipeline.permission: 0
pipeline.projectid: 0
pipeline.stages.len: 1
pipeline.stages.stages0.buildorder: 1
pipeline.stages.stages0.enabled: true
pipeline.stages.stages0.id: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.actions.len: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.enabled: false
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.final: false
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.id: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.lastmodified: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.name: Script
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.len: 1
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.id: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.name: script
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.type: text
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.parameters.parameters0.value: echo lol
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.requirements.len: 0
pipeline.stages.stages0.jobs.jobs0.action.actions.actions0.type: Builtin
pipeline.stages.stages0.jobs.jobs0.action.actions.len: 1
pipeline.stages.stages0.jobs.jobs0.action.description: This is job 1
pipeline.stages.stages0.jobs.jobs0.action.enabled: false
pipeline.stages.stages0.jobs.jobs0.action.final: false
pipeline.stages.stages0.jobs.jobs0.action.id: 0
pipeline.stages.stages0.jobs.jobs0.action.lastmodified: 0
pipeline.stages.stages0.jobs.jobs0.action.name: Job 1
pipeline.stages.stages0.jobs.jobs0.action.parameters.len: 0
pipeline.stages.stages0.jobs.jobs0.action.requirements.len: 0
pipeline.stages.stages0.jobs.jobs0.enabled: false
pipeline.stages.stages0.jobs.jobs0.lastmodified: 0
pipeline.stages.stages0.jobs.jobs0.pipelineactionid: 0
pipeline.stages.stages0.jobs.jobs0.pipelinestageid: 0
pipeline.stages.stages0.jobs.len: 1
pipeline.stages.stages0.lastmodified: 0
pipeline.stages.stages0.name: stage 1
pipeline.stages.stages0.pipelinebuildjobs.len: 0
pipeline.stages.stages0.pipelineid: 0
pipeline.stages.stages0.prerequisites.len: 0
pipeline.stages.stages0.runjobs.len: 0
pipeline.type: build
`
	assert.Equal(t, expected, out.String())
	assert.NoError(t, err)
}

type Result struct {
	Body     string      `json:"body,omitempty" yaml:"body,omitempty"`
	BodyJSON interface{} `json:"bodyjson,omitempty" yaml:"bodyjson,omitempty"`
}

func TestWeird(t *testing.T) {

	r := Result{}
	r.Body = "foo"
	r.BodyJSON = map[string]interface{}{
		"cardID":      "1234",
		"items":       []string{},
		"description": "yolo",
	}

	expected := `result.body: foo
result.bodyjson.cardid: 1234
result.bodyjson.description: yolo
result.bodyjson.items.len: 0
result.bodyjson.len: 3
`

	out := &bytes.Buffer{}
	err := dump.Fdump(out, r, dump.WithDefaultLowerCaseFormatter())
	assert.NoError(t, err)
	assert.Equal(t, expected, out.String())

	dump.Dump(r)
}
