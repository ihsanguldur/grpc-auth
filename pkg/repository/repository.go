package repository

import (
	"gorm.io/gorm"
	"grpc-auth/pkg/api/models"
)

type Storage interface {
	GetUser(user *models.User) error
	CreateUser(user models.User) error
}

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{db: db}
}

func (s *storage) GetUser(user *models.User) error {
	if err := s.db.Where(&models.User{UserName: user.UserName, Password: user.Password}).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (s *storage) CreateUser(user models.User) error {
	if err := s.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}
