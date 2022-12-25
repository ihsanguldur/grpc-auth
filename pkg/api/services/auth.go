package services

import (
	"errors"
	"grpc-auth/pkg/api/models"
)

type AuthService interface {
	Login(user *models.User) error
	Register(user models.User) error
	Validate(user *models.User) error
}

type AuthRepository interface {
	GetUser(user *models.User) error
	CreateUser(user models.User) error
}

type authService struct {
	storage AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{storage: repo}
}

func (a *authService) Login(user *models.User) error {

	if user.UserName == "" {
		return errors.New("server[auth]-Username required")
	} else if user.Password == "" {
		return errors.New("server[auth]-Password required")
	}

	if err := a.storage.GetUser(user); err != nil {
		return errors.New("server[auth]-user does not found")
	}

	return nil
}

func (a *authService) Register(user models.User) error {
	if user.UserName == "" {
		return errors.New("server[auth]-Username required")
	} else if user.Password == "" {
		return errors.New("server[auth]-Password required")
	}

	if err := a.storage.CreateUser(user); err != nil {
		return errors.New("server[auth]-error when creating new user")
	}

	return nil
}

func (a *authService) Validate(user *models.User) error {
	if user.UserName == "" {
		return errors.New("server[auth]-Username required")
	}

	if err := a.storage.GetUser(user); err != nil {
		return errors.New("server[auth]-user not found")
	}

	return nil
}
