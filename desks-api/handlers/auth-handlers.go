package handlers

import (
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util"
	"github.com/ChristophBe/desks/desks-api/util/configuration"
	"net/http"
)

type AuthHandlers interface {
	Redirect(writer http.ResponseWriter, request *http.Request)
	Callback(writer http.ResponseWriter, request *http.Request)
}

type authHandlersImpl struct {
	oicdService services.OICDService
}

func (a authHandlersImpl) Redirect(writer http.ResponseWriter, request *http.Request) {

	state := ""
	url := a.oicdService.GetRedirectUrl(state)

	http.Redirect(writer, request, url, http.StatusTemporaryRedirect)
}

func (a authHandlersImpl) Callback(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	//state := request.URL.Query().Get("state")

	// Verify state.

	idToken, err := a.oicdService.GetIdToken(ctx, request.URL.Query().Get("code"))

	if err != nil {
		err := util.Unauthorized(nil)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	// Extract custom claims.
	var claims struct {
		Email      string `json:"email"`
		Verified   bool   `json:"email_verified"`
		FamilyName string `json:"family_name"`
		GivenName  string `json:"given_name"`
	}
	if err := idToken.Claims(&claims); err != nil {
		err := util.Unauthorized(nil)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	userRepository := data.NewUserRepository()

	user, err := userRepository.FindByExternalUserID(request.Context(), idToken.Subject)

	if err != nil {
		user = models.User{}
	}

	user.ExternalUserId = idToken.Subject
	user.FamilyName = claims.FamilyName
	user.GivenName = claims.GivenName

	if user, err = userRepository.Save(request.Context(), user); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}

	if err = services.NewAuthCookieService().SetAuthCookieFor(user, writer, request); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
		return
	}

	baseUrl := configuration.BaseUrl.GetValue()

	http.Redirect(writer, request, baseUrl, http.StatusTemporaryRedirect)
}

func NewAuthHandlers(oicdService services.OICDService) AuthHandlers {
	return &authHandlersImpl{oicdService: oicdService}
}
