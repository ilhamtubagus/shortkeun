package persistence

import (
	"github.com/ilhamtubagus/urlShortener/domain/entity"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mgm.Collection
}

func (c userRepository) SaveUser(user *entity.User) error {
	err := c.collection.Create(user)
	if err != nil {
		return err
	}
	return nil
}
func (c userRepository) UpdateUser(user *entity.User) error {
	err := c.collection.Update(user)
	if err != nil {
		return err
	}
	return nil

}
func (c userRepository) FindUserByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	// err := c.collection.First(bson.M{operator.And: bson.A{bson.M{"email": email}, bson.M{"activation_code": bson.ErrDecodeToNil}}}, user)
	err := c.collection.First(bson.M{"email": email}, user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func NewUserRepository(c *mgm.Collection) userRepository {
	return userRepository{c}
}
