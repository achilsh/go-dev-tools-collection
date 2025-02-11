package main

import (
	"fmt"
	"log"
	"net/http"
)

// 处理带参数的 GET 请求
func paramHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "只支持 GET 请求", http.StatusMethodNotAllowed)
		return
	}
	// 获取查询参数
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "访客"
	}
	fmt.Println("call once..")
	fmt.Fprintf(w, "你好，%s！欢迎访问带参数的 GET 服务器。", name)
}

func main() {
	http.HandleFunc("/hello", paramHandler)
	addr := ":6666"
	fmt.Printf("服务器正在监听端口 %s...\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
