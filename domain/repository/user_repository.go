package repository

import "github.com/ilhamtubagus/urlShortener/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
}
