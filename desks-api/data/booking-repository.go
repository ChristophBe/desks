package data

import (
	"context"
	"github.com/ChristophBe/desks/desks-api/models"
	"gorm.io/gorm"
	"time"
)

type BookingRepository interface {
	Save(ctx context.Context, booking models.Booking) (models.Booking, error)
	FetchByUserId(ctx context.Context, userId int64) ([]models.Booking, error)
	FetchByRoomAndDate(ctx context.Context, roomId int64, date time.Time) ([]models.Booking, error)
	AnonymizeOldBookings(ctx context.Context, age int) (int64, error)
}

func NewBookingRepository() BookingRepository {
	return new(bookingRepositoryImpl)
}

type bookingRepositoryImpl struct{}

func (b bookingRepositoryImpl) AnonymizeOldBookings(ctx context.Context, maxAgeInDays int) (affectedRows int64, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	tx := db.Model(new(models.Booking)).Where("updated < NOW() - ? * INTERVAL '1 days'", maxAgeInDays).Update("user_id", gorm.Expr("NULL"))
	return tx.RowsAffected, tx.Error
}

func (b bookingRepositoryImpl) FetchByRoomAndDate(ctx context.Context, roomId int64, date time.Time) (bookings []models.Booking, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Preload("User").Preload("Room").Where("room_id = ? and start BETWEEN ?::timestamp and ?::timestamp ", roomId, date, date.Add(24*time.Hour)).Order("start").Find(&bookings).Error
	return
}
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
