package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// 定义要请求的 URL
	url := "http://localhost:6666/hello?name=aaaa"

	// 发送 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("请求出错: %v\n", err)
		return
	}
	// 确保在函数结束时关闭响应体
	defer resp.Body.Close()

	// 读取响应体内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return
	}

	// 打印响应状态码和响应体内容
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))
}
