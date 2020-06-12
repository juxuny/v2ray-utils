package cmd

import (
	"encoding/json"
	"strings"
)

func ToJson(v interface{}, pretty ...bool) string {
	if len(pretty) > 0 && pretty[0] {
		data, _ := json.MarshalIndent(v, "", "  ")
		return string(data)
	}
	data, _ := json.Marshal(v)
	return string(data)
}

// isIgnore is used to check if the server name contained a invalid keyword
func isIgnore(tag string, ignoreKeys []string) bool {
	for _, k := range ignoreKeys {
		if strings.Contains(tag, k) {
			return true
		}
	}
	return false
}
