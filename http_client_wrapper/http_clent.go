package httpclientwrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	logger "github.com/achilsh/go-dev-tools-collection/base-lib/log"
	"github.com/go-resty/resty/v2"
)

// HttpOptions http client 的请求options
type HttpOptions struct {
	Headers map[string]string
	// 如果 设置了 BaseUrl; 那么此处的url 是相对路径，比如： /api/v1/abc 整个路径是：BaseUrl + Url
	// 如果 没有设置 BaseUrl; 那么此处的url 是全部完整路径，比如： https://www.xxx.com/api/v1/abc
	Url     string
	IsDebug bool
	Timout  time.Duration // 请求回包超时控制
	BaseUrl string        //统一接口前缀，公共部分, 比如: https://www.xxx.com ； 设置该值的作用：不需要在每个请求中重复写主机名和公共前缀。
}

// Option 设置 http client 的请求 options的方法抽象
type Option func(*HttpOptions)

// WithHeaders 设置 http client 请求的http header字段
func WithHeaders(heads map[string]string) Option {
	return func(opts *HttpOptions) {
		if opts == nil {
			return
		}
		for k, v := range heads {
			opts.Headers[k] = v
		}
	}
}

// WithDebug 设置http client 请求 debug 选项
func WithDebug(enable bool) Option {
	return func(opts *HttpOptions) {
		if opts == nil {
			return
		}
		opts.IsDebug = enable
	}
}

// WithTimeOut 设置请求耗时控制
func WithTimeOut(tm time.Duration) Option {
	if tm <= 0.0 {
		tm = 40 * time.Millisecond
	}
	return func(opts *HttpOptions) {
		if opts == nil {
			return
		}
		opts.Timout = tm
	}
}

// WithBaseUrl 设置base url
func WithBaseUrl(baseUrl string) Option {
	return func(opts *HttpOptions) {
		if baseUrl == "" || opts == nil {
			return
		}
		opts.BaseUrl = baseUrl
	}
}

// WithUrl 设置 http client 请求 url 值
func WithUrl(url string) Option {
	return func(opts *HttpOptions) {
		if opts == nil {
			return
		}
		if url == "" {
			return
		}
		opts.Url = url
	}
}

// HttpClientRpcType 定义http client 模板函数类型
type HttpClientRpcType[I any, O any] func(ctx context.Context, reqBody *I, opts ...Option) (*O, error)

// HttpClientCallTpl http client request 调用， I 是 request body 类型, O 是 response body 类型
func HttpClientCallTpl[I any, O any](ctx context.Context, reqBody *I, opts ...Option) (*O, error) {
	defaultOpt := &HttpOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsDebug: false,
		Timout:  40 * time.Millisecond,
	}
	//
	for _, opt := range opts {
		opt(defaultOpt)
	}
	if defaultOpt.Url == "" {
		return nil, fmt.Errorf("url is nil")
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("format req data fail, err: %v", err)
		return nil, fmt.Errorf("format reqbody fail")
	}
	var cli *resty.Client
	if defaultOpt.BaseUrl != "" {
		cli = resty.New().SetBaseURL(defaultOpt.BaseUrl)
	} else {
		cli = resty.New()
	}
	r := cli.R().
		SetDebug(defaultOpt.IsDebug).
		SetContext(ctx).
		SetHeaders(defaultOpt.Headers).
		SetBody(reqBodyData).
		SetResult(new(O))
	var cancelFn context.CancelFunc
	if defaultOpt.Timout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, defaultOpt.Timout)
		r = r.SetContext(ctx)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	response, err := r.Post(defaultOpt.Url)
	if err != nil {
		logger.Errorf("post url: [%s] fail, err: %v", defaultOpt.Url, err)
		return nil, err
	}
	if response.RawResponse.StatusCode != http.StatusOK {
		logger.Errorf("err code: %v, msg: %v", response.RawResponse.StatusCode, response.RawResponse.Status)
		return nil, fmt.Errorf("error code: %v", response.RawResponse.StatusCode)
	}
	repItemPtr, ok := response.Result().(*O)
	if !ok {
		logger.Errorf("response type is not right")
		return nil, fmt.Errorf("response type is not right")
	} else if repItemPtr == nil {
		logger.Errorf("response is nil")
		return nil, fmt.Errorf("response is nil")
	}
	return repItemPtr, nil
}

func HttpClientGetTpl(ctx context.Context, queryParam map[string]string, opts ...Option) (int, string, error) {
	defaultOpt := &HttpOptions{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IsDebug: false,
		Timout:  100 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(defaultOpt)
	}

	if defaultOpt.Url == "" {
		return 400, "url is empty", fmt.Errorf("url is nil")
	}

	var cli *resty.Client
	if defaultOpt.BaseUrl != "" {
		cli = resty.New().SetBaseURL(defaultOpt.BaseUrl)
	} else {
		cli = resty.New()
	}

	r := cli.R().SetDebug(defaultOpt.IsDebug).SetContext(ctx).SetHeaders(defaultOpt.Headers)
	if len(queryParam) > 0 {
		r = r.SetQueryParams(queryParam)
	}
	var cancelFn context.CancelFunc
	if defaultOpt.Timout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, defaultOpt.Timout)
		r = r.SetContext(ctx)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	response, err := r.Get(defaultOpt.Url + "/")
	if err != nil {
		return 500, "get error", err
	}
	return response.StatusCode(), response.Status(), err
}
