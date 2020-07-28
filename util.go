package form3_sdk

import (
	"bytes"
	"fmt"
	"reflect"
)

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
