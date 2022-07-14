package authentication

type Token struct {
	// access token
	AccessToken string `json:"access_token"`
	// refresh token
	RefreshToken string `json:"refresh_token"`
}
