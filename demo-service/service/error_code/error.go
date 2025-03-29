package error_code

import (
	errors "demo-service/service/utils/error_def"
)

const (
	ErrCode_WelcomeWord = "5001"
)

var (
	WelcomeWordError = &errors.LogicError{
		Code:    ErrCode_WelcomeWord,
		Message: "load welcome word fail",
	}
)
