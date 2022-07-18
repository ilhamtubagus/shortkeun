package user

type UserService interface {
	FindUserByEmail(email string) (*User, error)
	SaveUser(user *User) error
	UpdateUser(user *User) error
	ActivateAccount(email, activationCode string) (*User, error)
}
type userService struct {
	userRepository UserRepository
}

func (service userService) SaveUser(user *User) error {
	return service.userRepository.SaveUser(user)
}

func (service userService) FindUserByEmail(email string) (*User, error) {
	return service.userRepository.FindUserByEmail(email)
}

func (service userService) UpdateUser(user *User) error {
	return service.userRepository.UpdateUser(user)
}

func (service userService) ActivateAccount(email, activationCode string) (*User, error) {
	user, err := service.userRepository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserService(userRepository UserRepository) UserService {
	return userService{
		userRepository: userRepository,
	}
}
