package handlers

import (
	"fmt"
	"github.com/ChristophBe/rooms/rooms-api/data"
	"github.com/ChristophBe/rooms/rooms-api/middlewares"
	"github.com/ChristophBe/rooms/rooms-api/models"
	"github.com/ChristophBe/rooms/rooms-api/util"
	"net/http"
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

	if room, err := roomRepository.FindById(ctx, body.Room.Id); err != nil || room.Id == 0 {
		err = util.BadRequest(err)
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

	if err = util.JsonResponseWriter(bookings, http.StatusAccepted, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}
}
