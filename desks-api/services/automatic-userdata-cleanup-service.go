package services

import (
	"context"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/util"
	"log"
	"time"
)

type AutomaticUserDataCleanupService interface {
	InitAutomaticUserDataCleaner(ctx context.Context)
}

func NewAutomaticUserDataCleanupService() AutomaticUserDataCleanupService {
	service := new(automaticUserDataCleanupServiceImpl)
	service.userRepository = data.NewUserRepository()
	service.bookingRepository = data.NewBookingRepository()
	return service
}

type automaticUserDataCleanupServiceImpl struct {
	userRepository    data.UserRepository
	bookingRepository data.BookingRepository
}

func (a automaticUserDataCleanupServiceImpl) InitAutomaticUserDataCleaner(ctx context.Context) {

	automaticCleanerInterval := util.GetIntEnvironmentVariable("USERDATA_CLEANER_INTERVAL_HOURS", 24)

	ticker := time.Tick(time.Duration(automaticCleanerInterval) * time.Hour)
	go a.cleanUserData(ctx)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				log.Println("clean old user data")
				a.cleanUserData(ctx)
			}
		}
	}()
}

func (a automaticUserDataCleanupServiceImpl) cleanUserData(ctx context.Context) {

	defer func() {
		recovered := recover()
		if recovered != nil {
			log.Println("Recovered panic in user data cleaner ", recovered)
		}
	}()
	maxDataAge := util.GetIntEnvironmentVariable("MAX_USERDATA_AGE_DAYS", 90)
	deletedUsers, err := a.userRepository.DeleteInactiveUser(ctx, maxDataAge)
	if err != nil {
		log.Println("failed to cleanup inactive users", err)
	}
	log.Printf("removed %d inactive users", deletedUsers)

	anonymizeBookings, err := a.bookingRepository.AnonymizeOldBookings(ctx, maxDataAge)
	if err != nil {
		log.Println("failed to cleanup old bookings", err)
	}
	log.Printf("anonymized %d old bookings", anonymizeBookings)
}
