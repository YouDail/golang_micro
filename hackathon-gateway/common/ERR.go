package common

import (
	"encoding/json"
	"strings"
)

//错误信息结构体
type ERR struct {
	Id     string `json:"id"`
	Code   int64  `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

//提取detail
func GetErr(s string) string {

	var d ERR
	err := json.Unmarshal([]byte(s), &d)
	if err != nil {

		if strings.Contains(s, `""`) {
			//return fmt.Printf("%q", s)

			strings.ReplaceAll(s, `""`, "`")
		}
		return s

	}

	if d.Detail == "not found" {
		return "内部错误"
	}

	if d.Detail == "sql: no rows in result set" {
		return "查询数据为空"
	}

	if d.Detail == "context deadline exceeded" {
		return "内部错误"
	}

	if strings.Contains(d.Detail, `""`) {
		strings.ReplaceAll(d.Detail, `""`, "`")
	}
	return d.Detail
}
