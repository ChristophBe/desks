package services

import (
	"fmt"
	"github.com/ChristophBe/rooms/rooms-api/models"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

const authCookieName = "auth_token"
const userIdClaimName = "user_id"

type AuthCookieService interface {
	SetAuthCookieFor(user models.User, writer http.ResponseWriter, request *http.Request) error
	GetUserIdFromCookie(request *http.Request) (int64, error)
}

func NewAuthCookieService() AuthCookieService {
	service := new(authCookieServiceImpl)
	service.jwtTokenService = NewJwtTokenService()
	return service
}

type authCookieServiceImpl struct {
	jwtTokenService JwtTokenService
}

func (a authCookieServiceImpl) SetAuthCookieFor(user models.User, writer http.ResponseWriter, request *http.Request) error {
	claims := jwt.MapClaims{}
	claims[userIdClaimName] = user.Id

	authToken, err := a.jwtTokenService.GenerateToken(claims, time.Hour*24)
	if err != nil {
		return fmt.Errorf("failed to generate auth token %w", err)
	}
	authCookie := http.Cookie{
		HttpOnly: true,
		Name:     authCookieName,
		Value:    authToken,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}
	request.AddCookie(&authCookie)
	http.SetCookie(writer, &authCookie)
	return nil
}

func (a authCookieServiceImpl) GetUserIdFromCookie(request *http.Request) (userId int64, err error) {
	authCookie, err := request.Cookie(authCookieName)
	if err != nil {
		err = fmt.Errorf("failed to read auth cookie %w", err)
		return
	}
	claims, err := a.jwtTokenService.VerifyToken(authCookie.Value)

	if err != nil {
		err = fmt.Errorf("failed to verify auth token %w", err)
		return
	}
	userIdClaim := claims[userIdClaimName]
	if userIdClaim == nil {
		err = fmt.Errorf("user id claim is empty")
		return
	}
	userId = int64(userIdClaim.(float64))
	return
}
