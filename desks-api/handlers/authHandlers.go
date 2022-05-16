package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/models"
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util"
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

func getOauthConfig() (oauthConfig *oauth2.Config, err error) {
	var (
		clientId         string
		clientSecret     string
		tokenUrl         string
		authorizationUrl string
		baseUrl          string
	)

	clientId, err = util.GetRequiredStringEnvironmentVariable("OAUTH_CLIENT_ID")
	if err != nil {
		return
	}
	clientSecret, err = util.GetRequiredStringEnvironmentVariable("OAUTH_CLIENT_SECRET")
	if err != nil {
		return
	}
	tokenUrl, err = util.GetRequiredStringEnvironmentVariable("OAUTH_TOKEN_URL")
	if err != nil {
		return
	}
	authorizationUrl, err = util.GetRequiredStringEnvironmentVariable("OAUTH_AUTHORIZATION_URL")
	if err != nil {
		return
	}
	baseUrl, err = util.GetRequiredStringEnvironmentVariable("BASE_URL")
	if err != nil {
		return
	}
	return &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid"},
		RedirectURL:  fmt.Sprintf("%s/auth/token", baseUrl),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authorizationUrl,
			TokenURL: tokenUrl,
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

	userInfoUrl, err := util.GetRequiredStringEnvironmentVariable("OAUTH_USERINFO_URL")
	if err != nil {
		err = util.InternalServerError(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}
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

	services.NewAuthCookieService().SetAuthCookieFor(user, writer, request)

	baseUrl, err := util.GetRequiredStringEnvironmentVariable("BASE_URL")
	if err != nil {
		err = util.InternalServerError(err)
		util.ErrorResponseWriter(err, writer, request)
		return
	}

	http.Redirect(writer, request, baseUrl, http.StatusTemporaryRedirect)
}
