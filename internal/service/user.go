package service

import (
	"errors"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/internal/model/entity"
	"filmLibraryVk/internal/repository"
	"filmLibraryVk/pkg"
	"log"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserById(id int) (presenter.UserResponse, error) {
	return s.repo.GetUserById(id)
}

func (s *UserService) GetUsers() ([]presenter.UserResponse, error) {
	return s.repo.GetUsers()
}

func (s *UserService) Login(login presenter.Login) (string, error) {
	user, err := s.repo.GetUserByUsername(login.Username)

	if err != nil {
		log.Printf("Invalid login or password")
		return "", errors.New("Invalid login or password")
	}

	err = pkg.ComparePasswords(user.Password, login.Password)

	if err != nil {
		log.Printf("Invalid login or password")
		return "", errors.New("Invalid login or password")
	}

	jwt, err := pkg.GenerateJWT(entity.User{
		Id:       user.Id,
		Username: user.Username,
		RoleId:   user.RoleId})
	if err != nil {
		log.Printf("Can not create JWT token for user %d", user.Id)
	}
	return jwt, nil
}

func (s *UserService) Register(register presenter.Register) (string, error) {
	var err error
	register.Password, err = pkg.EncodePassword(register.Password)

	log.Printf("Pass: %s", register.Password)
	if err != nil {
		log.Printf("Can not encode password")
		return "", nil
	}

	id, err := s.repo.CreateUser(register)
	if err != nil {
		return "", errors.New("User with such username already exists")
	}

	jwt, err := pkg.GenerateJWT(entity.User{
		Id:       id,
		Username: register.Username,
		RoleId:   1})
	if err != nil {
		log.Printf("Can not create JWT token for user %d", id)
		return "", err
	}
	return jwt, err
}

func (s *UserService) PutUser(id int, request presenter.UserRequest) (presenter.UserResponse, error) {
	pass, err := pkg.EncodePassword(*request.Password)
	request.Password = &pass

	if err != nil {
		return presenter.UserResponse{}, errors.New("Can not encode password")
	}

	return s.repo.PutUser(id, request)
}
func (s *UserService) PatchUser(id int, request presenter.UserRequest) (presenter.UserResponse, error) {
	if request.Password != nil {
		pass, err := pkg.EncodePassword(*request.Password)
		request.Password = &pass

		if err != nil {
			return presenter.UserResponse{}, errors.New("Can not encode password")
		}
	}
	return s.repo.PatchUser(id, request)
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}
