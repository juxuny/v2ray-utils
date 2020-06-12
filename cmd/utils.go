package cmd

import "strings"

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
