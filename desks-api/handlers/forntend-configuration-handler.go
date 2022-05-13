package handlers

import (
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util"
	"net/http"
)

func GetFrontendConfiguration(writer http.ResponseWriter, response *http.Request) {
	configuration := services.NewFrontendConfigurationService().GetFrontendConfiguration()
	if err := util.JsonResponseWriter(configuration, http.StatusOK, writer, response); err != nil {
		util.ErrorResponseWriter(util.InternalServerError(err), writer, response)
	}
}
