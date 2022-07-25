package handlers

import (
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/util"
	"net/http"
)

func GetUsersMe(user models.User, writer http.ResponseWriter, request *http.Request) {
	if err := util.JsonResponseWriter(user, http.StatusAccepted, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}
}
