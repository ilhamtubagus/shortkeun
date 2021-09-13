package lib

import (
	"testing"

	"github.com/ilhamtubagus/urlShortener/email"
	"github.com/stretchr/testify/assert"
)

func TestSendHTMLMail(t *testing.T) {
	LoadEnv("../.env")
	body := email.RegistrationMailBody{UserAgent: "Firefox on Desktop", IP: "12.23.12", DateTime: "Friday, Sep 3, 2021 9:27:56 AM (WIB)", Code: "12332"}
	err := SendHTMLMail([]string{"ilhamta@gmail.com"}, "Registration", body, "../email/template/registrationMail.html", []string{"../public/companylogo.png"})
	assert.Equal(t, nil, err)
}
