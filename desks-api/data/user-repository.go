package data

import (
	"context"
	"github.com/ChristophBe/desks/desks-api/models"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (models.User, error)
	Save(ctx context.Context, user models.User) (models.User, error)
	FindById(ctx context.Context, id int64) (models.User, error)
	DeleteInactiveUser(ctx context.Context, maxAgeInDays int) (int64, error)
}

func NewUserRepository() UserRepository {
	return new(userRepositoryImpl)
}

type userRepositoryImpl struct{}

func (u userRepositoryImpl) DeleteInactiveUser(ctx context.Context, maxAgeInDays int) (removedUsers int64, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	tx := db.Where("updated < NOW() - ? * INTERVAL '1 days'", maxAgeInDays).Delete(new(models.User))
	return tx.RowsAffected, tx.Error
}

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
