package handlers

import (
	"github.com/ChristophBe/desks/desks-api/middlewares"
	"github.com/ChristophBe/desks/desks-api/util"
	"log"
	"net/http"
)

func GetUsersMe(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(middlewares.UserContextKey)
	if err := util.JsonResponseWriter(user, http.StatusAccepted, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}
}

func PostUsersLogin(writer http.ResponseWriter, request *http.Request) {
	log.Println("Post body to login")

}
