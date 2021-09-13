package repositories

import (
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositories struct {
	collection *mgm.Collection
}

func (c UserRepositories) AddNewUser(registrant *entities.User) (*entities.User, error) {
	err := c.collection.Create(registrant)
	if err != nil {
		return nil, err
	}
	return registrant, nil
}
func (c UserRepositories) FindUserByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	err := c.collection.First(bson.M{"email": email}, user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func NewUserRepository(c *mgm.Collection) UserRepositories {
	return UserRepositories{c}
}
