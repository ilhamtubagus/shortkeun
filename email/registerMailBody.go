package email

type RegistrationMailBody struct {
	UserAgent string
	IP        string
	DateTime  string
	Code      string
	ExpireAt  string
}
