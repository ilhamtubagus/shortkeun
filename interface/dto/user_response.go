package dto

//	A default response with message that describe the response result
//	swagger:response defaultResponse
type UserResponse struct {
	// in: body
	Body UserResponseBody
}
type UserResponseBody struct {
	ID     string `json:"id,omitempty"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:",omitempty"`
	Role   string `json:"role"`
}
