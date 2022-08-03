package dto

//	A default response with message that describe the response result
//	swagger:response defaultResponse
type DefaultResponse struct {
	// in: body
	Body DefaultResponseBody
}

type DefaultResponseBody struct {
	// The response message
	Message string `json:"message"`
	// Error code
	Code int32 `json:"code"`
}

func NewDefaultResponse(msg string, code int32) *DefaultResponseBody {
	return &DefaultResponseBody{Message: msg, Code: code}
}
