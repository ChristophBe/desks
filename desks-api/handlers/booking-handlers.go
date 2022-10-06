package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util"
	"net/http"
	"reflect"
	"time"
)

func PostBooking(user models.User, writer http.ResponseWriter, request *http.Request) {

	bookingRepository := data.NewBookingRepository()
	ctx := request.Context()

	body, err := readBookingFromRequest(request)
	if err != nil {
		util.ErrorResponseWriter(err, writer, request)
	}

	if err = validateBooking(ctx, user, body); err != nil {
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	booking, err := bookingRepository.Save(ctx, body)

	if err = util.JsonResponseWriter(booking, http.StatusAccepted, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}
}

func readBookingFromRequest(request *http.Request) (models.Booking, error) {
	var body models.Booking

	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		err = util.BadRequest(err)
		return models.Booking{}, nil
	}
	return body, nil
}

func fetchOriginalBookingForModification(ctx context.Context, user models.User, request *http.Request) (booking models.Booking, err error) {
	bookingRepository := data.NewBookingRepository()
	bookingId, err := util.GetIntegerUrlParameter(request, "id")

	if err != nil || reflect.ValueOf(bookingId).IsZero() {
		return
	}

	booking, err = bookingRepository.FetchById(ctx, bookingId)
	if err != nil || reflect.ValueOf(booking.Id).IsZero() {
		return
	}

	if booking.User.Id != user.Id {
		err = fmt.Errorf("booking with id %d not found for current user", bookingId)
		return
	}
	return
}

func PatchBooking(user models.User, writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	bookingRepository := data.NewBookingRepository()
	booking, err := fetchOriginalBookingForModification(ctx, user, request)

	if err != nil {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	body, err := readBookingFromRequest(request)
	if err != nil {
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if booking.End.Before(time.Now()) {
		err := util.BadRequest(fmt.Errorf("not allowed to modifiy past bookings"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}
	modifiedBooking := booking.Patch(body)

	if booking.Start.Before(time.Now()) && (!booking.Start.Equal(modifiedBooking.Start) || booking.RoomId != modifiedBooking.RoomId) {
		err := util.BadRequest(fmt.Errorf("not allowed to modifiy room or start for ongoing bookings"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if err = validateBooking(ctx, user, modifiedBooking); err != nil {
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	storedBooking, err := bookingRepository.Save(ctx, modifiedBooking)
	if err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}

	if err = util.JsonResponseWriter(storedBooking, http.StatusAccepted, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}

}

func validateChangesForOngoingBookings(originalBooking models.Booking, changes models.Booking) bool {
	return !changes.Start.IsZero() || originalBooking.Start.Equal(changes.Start)
}

func validateBooking(ctx context.Context, user models.User, booking models.Booking) error {

	roomRepository := data.NewRoomRepository()
	bookingRepository := data.NewBookingRepository()
	if booking.User.Id != user.Id {
		return util.Forbidden(fmt.Errorf("not able to create a booking for another user"))
	}

	if !booking.Start.Before(booking.End) {
		return util.BadRequest(fmt.Errorf("start should before end"))

	}

	maxBookingDate := services.NewFrontendConfigurationService().GetFrontendConfiguration().MaximalBookingDate

	if !booking.End.Before(maxBookingDate) {
		return util.BadRequest(fmt.Errorf("the booking should be befor %v", maxBookingDate))

	}

	if room, err := roomRepository.FindById(ctx, booking.Room.Id); err != nil || room.Id == 0 {
		return util.BadRequest(err)

	}

	if existsOverlap, err := bookingRepository.ExistsOverlappingBooking(ctx, booking); existsOverlap || err != nil {
		if err != nil {
			return util.InternalServerError(err)
		} else {
			return util.BadRequest(fmt.Errorf("not allowed to have overlapping bookings"))
		}
	}
	return nil
}

func DeleteBooking(user models.User, writer http.ResponseWriter, request *http.Request) {
	bookingRepository := data.NewBookingRepository()

	ctx := request.Context()

	bookingId, err := util.GetIntegerUrlParameter(request, "id")

	if err != nil || reflect.ValueOf(bookingId).IsZero() {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	booking, err := bookingRepository.FetchById(ctx, bookingId)
	if err != nil || reflect.ValueOf(booking.Id).IsZero() {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if booking.User.Id != user.Id {
		err = util.NotFound(fmt.Errorf("booking with id %d not found for current user", bookingId))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if booking.Start.Before(time.Now()) {
		err = util.BadRequest(fmt.Errorf("not allowed to remove current or past bookings"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	err = bookingRepository.Delete(ctx, booking)
	if err != nil {
		err = util.InternalServerError(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func GetBookingsByUser(user models.User, writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	userId, err := util.GetIntegerUrlParameter(request, "id")

	if err != nil || userId != user.Id {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}
	bookings, err := data.NewBookingRepository().FetchByUserId(ctx, userId)
	if err != nil {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if err = util.JsonResponseWriter(bookings, http.StatusOK, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}
}
func GetBookingsByRoomAndDate(_ models.User, writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	roomId, err := util.GetIntegerUrlParameter(request, "id")

	if err != nil {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	dateQueryParameter := request.URL.Query().Get("date")
	if dateQueryParameter == "" {
		err = util.NotFound(fmt.Errorf("dateQueryParameter not set"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	date, err := time.Parse("2006-01-02", dateQueryParameter)
	if err != nil {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	bookings, err := data.NewBookingRepository().FetchByRoomAndDate(ctx, roomId, date)
	if err != nil {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if err = util.JsonResponseWriter(bookings, http.StatusOK, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}
}
func GetBookingDefaults(user models.User, writer http.ResponseWriter, request *http.Request) {
	userId, err := util.GetIntegerUrlParameter(request, "id")

	if err != nil {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}
	if userId != user.Id {

		err = util.Unauthorized(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	defaults, err := services.NewDetermineDefaultsService().DetermineBookingDefaults(request.Context(), user)
	if err != nil {
		err = util.NotFound(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if err = util.JsonResponseWriter(defaults, http.StatusOK, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}
}
