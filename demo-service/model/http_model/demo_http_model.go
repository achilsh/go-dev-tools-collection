package http_model

type RequestParam struct {
	// RequestParam's id.
	// Required: true
	Id   int32  `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}
type ResponseParam struct {
	Result string `json:"result"`
}

type FormReqParam struct {
	NoceStr     string `json:"noce_str" form:"noce_str"`
	AccessToken string `json:"access_token" form:"access_token"`

	// key is fileName, value is file content; 多文件， files: 是文件表单名
	FileContentMap map[string][]byte `json:"files" form:"files"`
}

type FormResponse struct {
	XYZ string `json:"xyz"`
	ABC int    `json:"abc"`
}
