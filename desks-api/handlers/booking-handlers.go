package handlers

import (
	"fmt"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/middlewares"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util"
	"net/http"
	"time"
)

func PostBooking(writer http.ResponseWriter, request *http.Request) {

	bookingRepository := data.NewBookingRepository()
	roomRepository := data.NewRoomRepository()

	ctx := request.Context()
	var body models.Booking

	if err := util.ReadJsonBody(request, &body); err != nil {
		err = util.BadRequest(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	user, ok := ctx.Value(middlewares.UserContextKey).(models.User)
	if !ok {
		err := util.Unauthorized(fmt.Errorf("can not get user from request contenxt"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if body.User.Id != user.Id {
		err := util.Forbidden(fmt.Errorf("not able to create a booking for another user"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if !body.Start.Before(body.End) {
		err := util.BadRequest(fmt.Errorf("start should before end"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	maxBookingDate := services.NewFrontendConfigurationService().GetFrontendConfiguration().MaximalBookingDate

	if !body.End.Before(maxBookingDate) {
		err := util.BadRequest(fmt.Errorf("the booking should be befor %v", maxBookingDate))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if room, err := roomRepository.FindById(ctx, body.Room.Id); err != nil || room.Id == 0 {
		err = util.BadRequest(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if existsOverlap, err := bookingRepository.ExistsOverlappingBooking(ctx, body); existsOverlap || err != nil {
		if err != nil {
			err = util.InternalServerError(err)
		} else {
			err = util.BadRequest(fmt.Errorf("not allowed to have overlapping bookings"))
		}
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	booking, err := bookingRepository.Save(ctx, body)
	if body.User.Id != user.Id {
		err = util.InternalServerError(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	if err = util.JsonResponseWriter(booking, http.StatusAccepted, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}
}

func GetBookingsByUser(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	user, ok := ctx.Value(middlewares.UserContextKey).(models.User)

	if !ok {
		err := util.Unauthorized(fmt.Errorf("can not get user from request contenxt"))
		util.ErrorResponseWriter(err, writer, request)
		return
	}

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
func GetBookingsByRoomAndDate(writer http.ResponseWriter, request *http.Request) {

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
