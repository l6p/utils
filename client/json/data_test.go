package json

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const Json = `
{
	"f1": "v1",
	"f2": 0,
	"f3": [10, 20, 30],
	"f4": {
		"f41": "v2",
		"f42": ["v3", "v4", "v5"],
		"f43": []
	},
	"f5": 9.9
}
`

func TestGet1(t *testing.T) {
	assertEqual(t, D(Json).GetJson(""), Json)
	assert.Equal(t, D(Json).GetInt("f2"), 0)
	assert.Equal(t, D(Json).GetString("f1"), "v1")
	assert.Equal(t, D(Json).GetFloat("f5"), 9.9)
	assert.Equal(t, D(Json).GetFloat("f3[1]"), 20.0)
	assert.Equal(t, D(Json).GetString("f4.f41"), "v2")
	assert.Equal(t, D(Json).GetString("f4.f42[0]"), "v3")
}

func TestGet2(t *testing.T) {
	d1 := D(Json).GetData("f4")
	assert.Equal(t, d1.GetString("f41"), "v2")

	d2 := D(Json).GetData("f4.f42")
	assert.Equal(t, d2.GetString("[2]"), "v5")

	d3 := d1.GetData("f41")
	assert.Equal(t, d3.GetString(""), "v2")
}

func TestLoad(t *testing.T) {
	data := Load("data_test.json")
	assert.Equal(t, data.GetInt("f4.f43[1].f432"), 60)
}

func TestSet1(t *testing.T) {
	const Expected = `
	{
		"f1": "dummy",
		"f2": 100,
		"f3": 99.9,
		"f4": {
			"f41": {"a": 1},
			"f42": ["v3", "v4", "v5"],
			"f43": 10
		},
		"f5": 9.9
	}
	`
	value := D(Json).
		SetString("f1", "dummy").
		SetInt("f2", 100).
		SetFloat("f3", 99.9).
		SetJson("f4.f41", `{"a": 1}`).
		SetInt("f4.f43", 10).
		GetJson("")
	assertEqual(t, Expected, value)
}

func TestAppend1(t *testing.T) {
	value := D(`[]`).Append("", `{"dummy": 1}`).GetJson("")
	assertEqual(t, `[{"dummy": 1}]`, value)
}

func TestAppend2(t *testing.T) {
	const Expected = `
	{
	 	"f1": "v1",
	 	"f2": 0,
	 	"f3": [10, 20, 30, 40, 9.9],
		"f4": {
			"f41": "v2",
			"f42": ["v3", "v4", "v5", "v6"],
			"f43": [{"dummy": 1}]
		},
		"f5": 9.9
	}
	`
	value := D(Json).
		AppendInt("f3", 40).
		AppendFloat("f3", 9.9).
		AppendString("f4.f42", "v6").
		Append("f4.f43", `{"dummy": 1}`).
		GetJson("")
	assertEqual(t, Expected, value)
}

func TestFilter1(t *testing.T) {
	Json := `
	{
		"f4": [
			{
				"f41": "v2",
				"f42": ["v3", "v4", "v5"]
			},
			{
				"f41": "v3",
				"f42": ["v6", "v7", "v8"]
			}
		],
		"f5": 9.9
	}
	`
	value := D(Json).
		Filter("f4", func(data *Data) bool {
			return data.GetString("f41") == "v3"
		}).
		Filter("f4[0].f42", func(data *Data) bool {
			return data.GetString("") == "v6"
		}).
		GetString("f4[0].f42[0]")
	assert.Equal(t, "v6", value)
}

func TestMerge1(t *testing.T) {
	const Expected = `
	{
	 	"f1": "v1",
	 	"f2": 0,
	 	"f3": "dummy",
		"f4": {
			"f43": [10, 20]
		},
		"f5": 9.9,
		"f6": []
	}
	`
	value := D(Json).Merge("", `
	{
		"f3": "dummy",
		"f4": {
			"f43": [10, 20]
		},
		"f5": 9.9,
		"f6": []
	}
	`).GetJson("")
	assertEqual(t, Expected, value)
}

func TestMerge2(t *testing.T) {
	const Expected = `
	{
	 	"f1": "v1",
	 	"f2": 0,
	 	"f3": [10, 20, 30],
		"f4": {
			"f41": "v2",
			"f42": ["v3", "v4", "v5"],
			"f43": 10,
			"f44": "dummy"
		},
		"f5": 9.9
	}
	`
	value := D(Json).Merge("f4", `
	{
		"f43": 10,
		"f44": "dummy"
	}
	`).GetJson("")
	assertEqual(t, Expected, value)
}

func TestDelete1(t *testing.T) {
	value := D(Json).Delete("").GetJson("")
	assertEqual(t, `{}`, value)
}

func TestDelete2(t *testing.T) {
	const Expected = `
	{
	 	"f1": "v1",
	 	"f2": 0,
	 	"f3": [10, 20]
	}
	`
	value := D(Json).Delete("f3").Delete("f4").Delete("f5").GetJson("")
	assertEqual(t, Expected, value)
}

func TestMarshal(t *testing.T) {
	data := D(Json)

	serializedBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	var data2 *Data
	if err := json.Unmarshal(serializedBytes, &data2); err != nil {
		panic(err)
	}

	assertEqual(t, data.GetJson(""), data2.GetJson(""))
}

func assertEqual(t *testing.T, expected string, actual string) {
	expectedBuffer := new(bytes.Buffer)
	if err := json.Compact(expectedBuffer, []byte(expected)); err != nil {
		panic(err)
	}
	actualBuffer := new(bytes.Buffer)
	if err := json.Compact(actualBuffer, []byte(actual)); err != nil {
		panic(err)
	}
	assert.Equal(t, expectedBuffer.String(), actualBuffer.String())
}
