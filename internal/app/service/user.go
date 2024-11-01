package service

type IUserRepository interface {
	CreateUser(name string, login string, password string, balance int) error
}

type UserService struct {
	userRepository IUserRepository
}

func (s *UserService) CreateUser(name string, login string, password string, balance int) error {
	return s.userRepository.CreateUser(name, login, password, balance)
}

func NewUserService(userRepository IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}
