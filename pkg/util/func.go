package util

import "encoding/json"

// ToJSONString 将任意数据转换为JSON字符串，不返回错误（出错时返回空字符串）
func ToJSONString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
