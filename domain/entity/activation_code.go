package entity

import "time"

type ActivationCode struct {
	Code     string    `json:"code" bson:"code"`
	IssuedAt time.Time `json:"issued_at" bson:"issued_at"`
	ExpireAt time.Time `json:"expireAt" bson:"expire_at"`
}
