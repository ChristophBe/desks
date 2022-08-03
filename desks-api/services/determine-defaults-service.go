package services

import (
	"context"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/dtos"
	"github.com/ChristophBe/desks/desks-api/models"
	"sort"
	"time"
)

type DetermineDefaultsService interface {
	DetermineBookingDefaults(ctx context.Context, user models.User) (dtos.BookingDefaultsDto, error)
}

func NewDetermineDefaultsService() DetermineDefaultsService {
	service := new(determineDefaultsServiceImpl)
	service.bookingRepository = data.NewBookingRepository()
	return service
}

type determineDefaultsServiceImpl struct {
	bookingRepository data.BookingRepository
}

func (d determineDefaultsServiceImpl) DetermineBookingDefaults(ctx context.Context, user models.User) (defaults dtos.BookingDefaultsDto, err error) {
	bookings, err := d.bookingRepository.FetchByUserId(ctx, user.Id)
	if err != nil {
		return
	}

	defaults.Room, err = d.findMostBookedRoom(bookings)
	if err != nil {
		return
	}
	defaults.Start, err = d.findMedianTime(bookings, func(booking models.Booking) time.Time { return booking.Start })
	if err != nil {
		return
	}
	defaults.End, err = d.findMedianTime(bookings, func(booking models.Booking) time.Time { return booking.End })
	if err != nil {
		return
	}

	return

}

func (d determineDefaultsServiceImpl) findMostBookedRoom(bookings []models.Booking) (models.Room, error) {
	if len(bookings) <= 0 {
		return models.Room{}, fmt.Errorf("can not determine most booked room for empty bookings")
	}
	countByRoom := make(map[int64]int, len(bookings))
	roomByRoomId := make(map[int64]models.Room, len(bookings))

	for _, booking := range bookings {
		room := booking.Room
		countByRoom[room.Id]++
		roomByRoomId[room.Id] = room
	}

	maxCount := 0
	maxRoomId := int64(0)

	for roomId, count := range countByRoom {

		if count > maxCount {
			maxRoomId = roomId
			maxCount = count
		}
	}

	return roomByRoomId[maxRoomId], nil
}

func (d determineDefaultsServiceImpl) findMedianTime(bookings []models.Booking, getTimeValue func(booking models.Booking) time.Time) (resultTime time.Time, err error) {

	if len(bookings) <= 0 {
		err = fmt.Errorf("can not determine median time for empty bookings")
		return
	}
	timeOffsets := make([]int, len(bookings))
	for i, booking := range bookings {
		bookingTime := getTimeValue(booking)
		hours, minute, _ := bookingTime.Clock()

		timeOffsets[i] = hours*60 + minute
	}

	sort.Ints(timeOffsets) // sort the timeOffsets

	index := len(timeOffsets) / 2

	medianTimeOffset := timeOffsets[index]
	now := time.Now()
	resultTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	resultTime = resultTime.Add(time.Duration(medianTimeOffset) * time.Minute)
	resultTime = resultTime.Add(24 * time.Hour)
	resultTime = resultTime.Round(15 * time.Minute)
	return
}
