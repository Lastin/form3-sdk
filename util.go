package form3

import (
	"bytes"
	"fmt"
	"reflect"
)

// Function to build filter string from generic struct
// It uses json tag, as the objects are filtered by attribute with name identical to that provided in json sent back by the API
func buildFilter(i interface{}) string {
	if i == nil {
		return ""
	}
	var buf bytes.Buffer
	v := reflect.ValueOf(i)
	for j := 0; j < v.NumField(); j++ {
		if v.Field(j).IsValid() && !v.Field(j).IsZero() {
			jsonTagName := v.Type().Field(j).Tag.Get("json")
			if len(jsonTagName) > 0 {
				buf.WriteString(fmt.Sprintf("&filter[%s]=%v", jsonTagName, reflect.Indirect(v.Field(j)).Interface()))
			}

		}
	}
	return buf.String()
}
