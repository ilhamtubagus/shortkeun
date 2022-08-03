package dto

//	A default response with message that describe the response result
//	swagger:response defaultResponse
type RegistrationResponse struct {
	// in: body
	Body RegistrationResponseBody
}
type RegistrationResponseBody struct {
	ID     string `json:"id,omitempty"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:",omitempty"`
	Role   string `json:"role"`
}
