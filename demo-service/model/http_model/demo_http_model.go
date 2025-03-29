package http_model

type RequestParam struct {
	Id   int32  `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}
type ResponseParam struct {
	Result string `json:"result"`
}
