package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logc"
)

// 自定义的中间件
func HeaderMiddlewares(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xabcHeader := r.Header.Get("x-abc")
		if len(xabcHeader) > 0 {
			xabcHeader += "______1123_abc"
		}
		r.Header.Set("x-abc", xabcHeader)
		r.Header.Add("X-Middleware", "static-middleware")

		r, _ = ModifyRequestBody(r)
		next(w, r)
	}
}

// 修改请求Body的函数（示例：给JSON请求添加字段）
func ModifyRequestBody(r *http.Request) (*http.Request, error) {
	// 仅处理POST/PUT等有Body的请求
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		return r, nil
	}

	// 读取原始Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("读取请求体失败: %v", err)
	}
	defer r.Body.Close()

	// 打印原始Body
	logc.Infof(r.Context(), "原始请求体: %s", string(body))

	// 解析JSON并添加字段（假设Body是JSON格式）
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// 如果不是JSON格式，不修改直接返回
		fmt.Printf("请求体不是JSON格式，不修改: %v\n", err)
		return r, nil
	}

	// 添加自定义字段（例如：添加网关标识）
	data["name"] = "data++++++1"
	// data["timestamp"] = time.Now().Unix()

	// 重新序列化JSON
	newBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("序列化修改后的Body失败: %v", err)
	}

	// 打印修改后的Body
	logc.Infof(r.Context(), "修改后的请求体: %s", string(newBody))

	// 重新构造请求对象（因为Body已被读取，需要重置）
	r.Body = io.NopCloser(bytes.NewBuffer(newBody))
	r.ContentLength = int64(len(newBody))
	r.Header.Set("Content-Length", fmt.Sprintf("%d", len(newBody)))

	return r, nil
}
