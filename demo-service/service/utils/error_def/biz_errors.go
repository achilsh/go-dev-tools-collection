package error_def

// CliErrorEr 定义错误码信息接口
type CliErrorEr interface {
	GetCode() string
	GetCodeMsg() string
	Error() string
	GetCodeArgs() []any
}

// HttpResponse represents response item.
// swagger:model HttpResponse
type HttpResponse struct {
	// Custom error code
	// required: true
	// example:
	ErrorCode string `json:"code"`
	// Custom error message description
	// required: true
	// example: "fail"
	ErrorMessage string `json:"msg"`
	// Custom request tracerId.
	// required: true
	// example: abc123iefdfa
	SeqId string `json:"seqId"`
	// Custom request cost time in millisecond.
	// required: true
	// example: 12
	CostTimeMs int64 `json:"costTimeMs"`
	// Custom request response context.
	// required: false
	// example: null
	Data any `json:"data,omitempty"`
}
