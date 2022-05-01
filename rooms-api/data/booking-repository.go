package data

import (
	"context"
	"github.com/ChristophBe/rooms/rooms-api/models"
)

type BookingRepository interface {
	Save(ctx context.Context, booking models.Booking) (models.Booking, error)
	FetchByUserId(ctx context.Context, userId int64) ([]models.Booking, error)
}

func NewBookingRepository() BookingRepository {
	return new(bookingRepositoryImpl)
}

type bookingRepositoryImpl struct{}

func (b bookingRepositoryImpl) FetchByUserId(ctx context.Context, userId int64) (bookings []models.Booking, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Preload("User").Preload("Room").Where("user_id = ?", userId).Order("start").Find(&bookings).Error
	return
}

func (b bookingRepositoryImpl) Save(ctx context.Context, booking models.Booking) (stored models.Booking, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Save(&booking).Error
	return booking, err
}
