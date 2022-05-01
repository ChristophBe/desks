package handlers

import (
	"fmt"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/middlewares"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/services"
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

	userRepository := data.NewUserRepository()
	var body models.User
	if err := util.ReadJsonBody(request, &body); err != nil {
		util.ErrorResponseWriter(util.BadRequest(err), writer, request)
		return
	}

	if len(body.Firstname) <= 0 && len(body.Firstname) <= 0 && len(body.Lastname) <= 0 {
		err := fmt.Errorf("invalid body")
		util.ErrorResponseWriter(util.BadRequest(err), writer, request)
		return
	}

	user, err := userRepository.FindByUsername(request.Context(), body.Username)

	if err != nil {
		user = models.User{}
	}

	user.Firstname = body.Firstname
	user.Lastname = body.Lastname
	user.Username = body.Username

	if user, err = userRepository.Save(request.Context(), user); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}

	services.NewAuthCookieService().SetAuthCookieFor(user, writer, request)
	if err = util.JsonResponseWriter(user, http.StatusAccepted, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}
}
