package handlers

import (
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/util"
	"net/http"
)

func GetAllRooms(writer http.ResponseWriter, request *http.Request) {
	roomsRepositories := data.NewRoomRepository()

	rooms, err := roomsRepositories.FindAll(request.Context())
	if err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}

	if err = util.JsonResponseWriter(rooms, http.StatusOK, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
	}
}
