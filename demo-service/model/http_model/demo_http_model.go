package http_model

type RequestParam struct {
	Id   int32  `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}
type ResponseParam struct {
	Result string `json:"result"`
}

type FormReqParam struct {
	NoceStr     string `json:"noce_str" form:"noce_str"`
	AccessToken string `json:"access_token" form:"access_token"`
	FileContent []byte `json:"file" form:"file"`

	// 文件的文件名，从表单字段解析, 内部参数，前端忽略
	FileName string `json:"-"`
}

type FormResponse struct {
	XYZ string `json:"xyz"`
	ABC int    `json:"abc"`
}
