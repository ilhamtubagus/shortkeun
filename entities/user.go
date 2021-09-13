package entities

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string
}

type GoogleCredential struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}
