package middlewares

import (
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util"
	"net/http"
)

type AuthorizedHandler func(user models.User, writer http.ResponseWriter, request *http.Request)

func AuthMiddleware(handler AuthorizedHandler) http.Handler {

	authCookieService := services.NewAuthCookieService()
	userRepository := data.NewUserRepository()

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		userId, err := authCookieService.GetUserIdFromCookie(request)
		if err != nil {
			err = util.Unauthorized(err)
			util.ErrorResponseWriter(err, writer, request)
			return
		}

		user, err := userRepository.FindById(ctx, userId)
		if err != nil || user.Id == 0 {
			err = util.Unauthorized(err)
			util.ErrorResponseWriter(err, writer, request)
			return
		}

		handler(user, writer, request.WithContext(ctx))
	})
}
