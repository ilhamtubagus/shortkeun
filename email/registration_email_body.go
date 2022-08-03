package email

type RegistrationEmailBody struct {
	UserAgent string
	IP        string
	DateTime  string
	Code      string
	ExpireAt  string
}
