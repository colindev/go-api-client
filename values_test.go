package api

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {

	jsonBytes := []byte(`{
	
		"a": {"a1": {"a2": 1}},
		"b": [{"b2": 2}],
		"c": 3

	}`)

	var d map[string]interface{}

	json.Unmarshal(jsonBytes, &d)

	v := &Values{d}

	asserts := map[string]interface{}{
		"a.a1.a2": 1,
		"b.0.b2":  2,
	}

	for key, val := range asserts {
		if x := v.Get(key); x != val {
			t.Error(
				"v.Get("+key+") = ", x, fmt.Sprintf("%T", x),
				"Expect", val, fmt.Sprintf("%T", val))
		}
	}
}
