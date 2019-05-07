package common

import "net/http"

// 统一错误输出接口
func ErrorHandle(err error, w http.ResponseWriter) {
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
