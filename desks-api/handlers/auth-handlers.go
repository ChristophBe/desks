package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util"
	"github.com/ChristophBe/desks/desks-api/util/configuration"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

type OpenIdUserInfo struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
}

type OauthState struct {
	RedirectUrl string `json:"redirectUrl"`
}

func getOauthConfig() (oauthConfig *oauth2.Config, err error) {
	return &oauth2.Config{
		ClientID:     configuration.OauthClientId.GetValue(),
		ClientSecret: configuration.OauthClientSecret.GetValue(),
		Scopes:       []string{"openid"},
		RedirectURL:  fmt.Sprintf("%s/auth/token", configuration.BaseUrl.GetValue()),
		Endpoint: oauth2.Endpoint{
			AuthURL:  configuration.OauthAuthorizationUrl.GetValue(),
			TokenURL: configuration.OauthTokenUrl.GetValue(),
		},
	}, nil

}
func AuthLogin(writer http.ResponseWriter, request *http.Request) {

	conf, err := getOauthConfig()
	if err != nil {
		err = util.InternalServerError(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	http.Redirect(writer, request, url, http.StatusTemporaryRedirect)

}

func AuthRedirect(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	conf, err := getOauthConfig()
	if err != nil {
		err = util.InternalServerError(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	code := request.URL.Query().Get("code")
	if code == "" {
		err := util.Unauthorized(nil)
		util.ErrorResponseWriter(err, writer, request)
		return
	}
	tok, err := conf.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {

		err = util.Unauthorized(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	userInfoUrl := configuration.OauthUserinfoUrl.GetValue()

	client := conf.Client(ctx, tok)
	resp, err := client.Get(userInfoUrl)
	if err != nil {
		err = util.Unauthorized(err)
		util.ErrorResponseWriter(err, writer, request)
		return

	}
	defer request.Body.Close()

	userInfoResponseBoy, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		err = util.Unauthorized(err)
		util.ErrorResponseWriter(err, writer, request)
	}

	var userInfo OpenIdUserInfo
	err = json.Unmarshal(userInfoResponseBoy, &userInfo)
	if err != nil {
		err = util.Unauthorized(err)
		util.ErrorResponseWriter(err, writer, request)
	}

	userRepository := data.NewUserRepository()

	user, err := userRepository.FindByExternalUserID(request.Context(), userInfo.Sub)

	if err != nil {
		user = models.User{}
	}

	user.ExternalUserId = userInfo.Sub
	user.FamilyName = userInfo.FamilyName
	user.GivenName = userInfo.GivenName

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
