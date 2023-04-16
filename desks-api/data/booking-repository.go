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
	FetchByRoomAndRange(ctx context.Context, id int64, from time.Time, to time.Time) ([]models.Booking, error)
	AnonymizeOldBookings(ctx context.Context, age int) (int64, error)
	FetchById(ctx context.Context, id int64) (models.Booking, error)
	Delete(ctx context.Context, booking models.Booking) error
	ExistsOverlappingBooking(ctx context.Context, booking models.Booking) (bool, error)
}

func NewBookingRepository() BookingRepository {
	return new(bookingRepositoryImpl)
}

type bookingRepositoryImpl struct{}

func (b bookingRepositoryImpl) ExistsOverlappingBooking(ctx context.Context, booking models.Booking) (result bool, err error) {

	idsToIgnore := []int64{booking.Id}

	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Raw("SELECT EXISTS(SELECT 1 FROM bookings WHERE room_id = ? and user_id = ? and ((\"start\" BETWEEN ?::timestamp and ?::timestamp) or (\"end\" BETWEEN ?::timestamp and ?::timestamp)) and id not in ?)", booking.Room.Id, booking.User.Id, booking.Start, booking.End, booking.Start, booking.End, idsToIgnore).Scan(&result).Error
	return
}

func (b bookingRepositoryImpl) Delete(ctx context.Context, booking models.Booking) (err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Delete(&booking).Error
	return
}

func (b bookingRepositoryImpl) FetchById(ctx context.Context, id int64) (booking models.Booking, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Preload("User").Preload("Room").Find(&booking, id).Error
	return
}

func (b bookingRepositoryImpl) AnonymizeOldBookings(ctx context.Context, maxAgeInDays int) (affectedRows int64, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	tx := db.Model(new(models.Booking)).Where("updated < NOW() - ? * INTERVAL '1 days'", maxAgeInDays).Update("user_id", gorm.Expr("NULL"))
	return tx.RowsAffected, tx.Error
}

func (b bookingRepositoryImpl) FetchByRoomAndRange(ctx context.Context, roomId int64, from, to time.Time) (bookings []models.Booking, err error) {
	db, err := GetConnection(ctx)
	if err != nil {
		return
	}
	err = db.Preload("User").Preload("Room").Where("room_id = ? and start BETWEEN ?::timestamp and ?::timestamp ", roomId, from, to).Order("start").Find(&bookings).Error
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
