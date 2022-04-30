package middlewares

import (
	"context"
	"github.com/ChristophBe/rooms/rooms-api/data"
	"github.com/ChristophBe/rooms/rooms-api/services"
	"github.com/ChristophBe/rooms/rooms-api/util"
	"net/http"
)

type userContextKeyType string

const UserContextKey = userContextKeyType("user")

func AuthMiddleware(handler http.Handler) http.Handler {

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

		ctx = context.WithValue(ctx, UserContextKey, user)
		handler.ServeHTTP(writer, request.WithContext(ctx))
	})
}
