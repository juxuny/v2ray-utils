package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 过滤掉非数字
func NumberFilter(s string) string {
	var data = []byte(s)
	ret := ""
	for _, b := range data {
		if b >= '0' && b <= '9' || b == '.' {
			ret += string([]byte{b})
		}
	}
	return ret
}

func getNetworkTraffic(s string) string {
	if len(s) == 0 {
		return ""
	}
	ret := ""
	index := strings.Index(s, "GB")
	for i := index - 1; i > 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			ret = string([]byte{s[i]}) + ret
		} else {
			break
		}
	}
	return ret
}

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

var (
	verbose bool
)

func logf(f string, _v ...interface{}) {
	if verbose {
		fmt.Printf(f+"\n", _v...)
	}
}

func log(_v ...interface{}) {
	if verbose {
		fmt.Println(_v...)
	}
}
