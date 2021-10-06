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
}

func NewDefaultResponse(msg string) *DefaultResponseBody {
	return &DefaultResponseBody{Message: msg}
}
