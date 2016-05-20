package bach

import (
	"testing"
)

func TestFlattenSimpleMap(t *testing.T) {
	d := make(map[string]interface{})
	d["foo"] = "bar"

	result := make(map[string]string)
	for e := range flattenMap(d, "") {
		result[e.Name] = e.Value
	}

	if result["foo"] != d["foo"] {
		t.Errorf("'%v' should have been '%v'", result["foo"], d["foo"])
	}
}

func TestFlattenNestedMap(t *testing.T) {
	d, err := LoadYaml("core/nested.yml")
	if err != nil {
		t.Fail()
	}

	result := make(map[string]string)
	for e := range flattenMap(d, "") {
		result[e.Name] = e.Value
	}

	t.Log(result)
	if result["FOO_BAR_BAZ"] != "hello world" {
		t.Errorf("'%v' should have been 'hello world'", result["foo"])
	}
}
