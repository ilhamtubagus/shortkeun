package entities

import "github.com/kamva/mgm/v3"

type Url struct {
	mgm.DefaultModel `bson:",inline"`
	Name             string   `bson:"name"`
	RedirectedUrl    string   `bson:"redirected_url"`
	ShortenedUrl     string   `bson:"shortened_url"`
	Tag              []string `bson:"tag"`
	ClickHistories   []Click  `bson:"click_histories"`
}

type Click struct {
	Location         string `bson:"location"`
	Ip               string `bson:"ip"`
	mgm.DefaultModel `bson:",inline"`
}
