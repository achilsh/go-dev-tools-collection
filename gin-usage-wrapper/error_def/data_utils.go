package error_def

import "encoding/json"

func StructToJsonString(item interface{}) string {
	if item == nil {
		return ""
	}
	buf, err := json.Marshal(item)
	if err != nil {
		return ""
	}
	return string(buf)
}
