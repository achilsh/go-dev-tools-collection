package error_def

import "fmt"

var _ CliErrorEr = (*LogicError)(nil)

// LogicError 返回给客户端的应答形式。 Status=1000 成功
type LogicError struct {
	Message string `json:"msg"`  // 错误消息描述
	Code    string `json:"code"` // 错误码
	Args    []any  `json:"-"`
}

func (e *LogicError) Error() string {
	if e == nil {
		return "un_init error obj"
	}
	return fmt.Sprintf("err_code: %v, err_msg: %v", e.Code, e.Message)
}
func (e *LogicError) GetCode() string {
	if e == nil {
		return "5000"
	}
	return e.Code
}
func (e *LogicError) GetCodeMsg() string {
	if e == nil {
		return "un defined"
	}
	return e.Message
}
func (e *LogicError) GetCodeArgs() []any {
	if e == nil {
		return nil
	}
	return e.Args
}

func buildErr(code string, msg string) {
	collectErrors.ErrorCollect[code] = &LogicError{
		Code:    code,
		Message: msg,
	}
}

var collectErrors = LogicErrors{
	ErrorCollect: make(map[string]CliErrorEr),
}

type LogicErrors struct {
	ErrorCollect map[string]CliErrorEr
}

// GetError 根据错误和可能的错误参数返回一个错误对象
func GetError(code string, errMsg ...any) CliErrorEr {
	formatStr := ""
	for _ = range errMsg {
		formatStr += "%v"
	}
	errItem, ok := collectErrors.ErrorCollect[code]
	if ok {
		if len(errMsg) <= 0 {
			return errItem
		}
		return &LogicError{
			Code:    code,
			Message: fmt.Sprintf(formatStr, errMsg...),
		}
	}
	if len(errMsg) <= 0 {
		return &LogicError{
			Code:    code,
			Message: fmt.Sprintf("code err: %v", code),
		}
	}
	composeErrMsg := fmt.Sprintf(formatStr, errMsg...)
	return &LogicError{
		Code:    code,
		Message: composeErrMsg,
	}
}

// GetErrorWithArgs 根据错误和可能的错误参数返回一个错误对象
func GetErrorWithArgs(code string, msgArgs ...any) CliErrorEr {
	return &LogicError{
		Code: code,
		Args: msgArgs,
	}
}
