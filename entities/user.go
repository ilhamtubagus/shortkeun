package entities

import (
	"github.com/kamva/mgm/v3"
)

type User struct {
	mgm.DefaultModel  `bson:",inline"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	Password          string `json:"password,omitempty"`
	Status            string `json:",omitempty"`
	*GoogleCredential `json:",omitempty" bson:",omitempty"`
}

type GoogleCredential struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}
