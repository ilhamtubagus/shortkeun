package dto

// Request schema for sign in with google account
// swagger:parameters googleSignIn
type SignInRequestGoogle struct {
	//
	// in: body
	// required: true
	Body GoogleSignInRequestBody
}

// swagger:model
type GoogleSignInRequestBody struct {
	// contain JWT ID Token obtained from google
	Credential string `json:"credential" validate:"required"`
}
