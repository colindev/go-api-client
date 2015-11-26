package api

import "strings"

type Values struct {
	body interface{}
}

func (v *Values) Get(key string) (ret interface{}) {
	indexes := strings.Split(key, ".")
	ret = v.body

	if key == "" {
		return
	}

	for _, ind := range indexes {
		ret = get(ret, ind)
		if ret == nil {
			return
		}
	}

	return
}

func get(d interface{}, key string) interface{} {

	v, ok := d.(map[string]interface{})
	if ok {
		return v[key]
	}

	return nil
}
