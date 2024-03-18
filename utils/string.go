package utils

import "strings"

// 字符串拼接
func JoinStrings(strSlice ...string) string {
	builder := strings.Builder{}
	for _, str := range strSlice {
		_, err := builder.WriteString(str)
		if err != nil {
			LogError("join string", err)
		}
	}
	return builder.String()
}
