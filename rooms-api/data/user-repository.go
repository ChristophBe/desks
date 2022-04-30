package data

import (
	"context"
	"github.com/ChristophBe/rooms/rooms-api/models"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (models.User, error)
	Save(ctx context.Context, user models.User) (models.User, error)
	FindById(ctx context.Context, id int64) (models.User, error)
}

func NewUserRepository() UserRepository {
	return new(userRepositoryImpl)
}

type userRepositoryImpl struct{}

func (u userRepositoryImpl) Save(ctx context.Context, user models.User) (stored models.User, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Save(&user).Error
	stored = user
	return
}

func (u userRepositoryImpl) FindByUsername(ctx context.Context, username string) (user models.User, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Where("username = ?", username).First(&user).Error
	return
}
func (u userRepositoryImpl) FindById(ctx context.Context, userId int64) (user models.User, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Where("id = ?", userId).First(&user).Error
	return
}
