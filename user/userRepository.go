package user

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(user *User) error
	UpdateUser(user *User) error
	FindUserByEmail(email string) (*User, error)
}
type userRepository struct {
	collection *mgm.Collection
}

func (c userRepository) CreateUser(user *User) error {
	err := c.collection.Create(user)
	if err != nil {
		return err
	}
	return nil
}
func (c userRepository) UpdateUser(user *User) error {
	err := c.collection.Update(user)
	if err != nil {
		return err
	}
	return nil

}
func (c userRepository) FindUserByEmail(email string) (*User, error) {
	user := &User{}
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
