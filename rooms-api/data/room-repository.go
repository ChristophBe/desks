package data

import (
	"context"
	"github.com/ChristophBe/rooms/rooms-api/models"
)

type RoomRepository interface {
	FindAll(ctx context.Context) ([]models.Room, error)
	FindById(ctx context.Context, id int64) (models.Room, error)
}

func NewRoomRepository() RoomRepository {
	return new(roomRepositoryImpl)
}

type roomRepositoryImpl struct{}

func (r roomRepositoryImpl) FindById(ctx context.Context, id int64) (room models.Room, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Find(&room, id).Error
	return
}

func (r roomRepositoryImpl) FindAll(ctx context.Context) (rooms []models.Room, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Find(&rooms).Error
	return
}
