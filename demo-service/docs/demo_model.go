package docs

import (
	"github.com/achilsh/go-dev-tools-collection/demo-service/model/http_model"
)

// swagger:route POST /demo-server/v1/x2 xxx v1x2DemoRequest
// call v1 x2 requestions..
// responses:
// 200: v1x2Response
// default: errResponse
// swagger:parameters v1x2DemoRequest
type v1x2ParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body http_model.RequestParam
}

// This text will appear as description of your response body.
// swagger:response v1x2Response
type v1x2ResponseWrapper struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`

	// in:body
	Body http_model.ResponseParam `json:"data"`
}

// This text will appear as description of your error response body.
// swagger:response errResponse
type errResponseWrapper struct {
	// Error code.
	Code int `json:"code"`
	// Error message.
	Message string `json:"message"`
}
