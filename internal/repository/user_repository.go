package repository

import (
	"gofiber-boilerplate/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Logger *logrus.Logger
}

func NewUserRepository(logger *logrus.Logger) *UserRepository {
	return &UserRepository{
		Logger: logger,
	}
}

func (r *UserRepository) CountByEmail(db *gorm.DB, email string) (int64, error) {
	var count int64
	err := db.Model(new(entity.User)).
		Where("email = ?", email).
		Count(&count).Error
	return count, err
}

func (r *UserRepository) FindByEmail(db *gorm.DB, email string) (*entity.User, error) {
	var user entity.User
	err := db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
