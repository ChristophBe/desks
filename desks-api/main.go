package main

import (
	"context"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/data"
	"github.com/ChristophBe/desks/desks-api/handlers"
	"github.com/ChristophBe/desks/desks-api/middlewares"
	"github.com/ChristophBe/desks/desks-api/services"
	"github.com/ChristophBe/desks/desks-api/util/configuration"
	"log"
	"net/http"
)
import "github.com/gorilla/mux"

func main() {
	initDefaultConfigurations()

	if err := data.InitDatabase(); err != nil {
		log.Fatal("Failed initialize Database", err)
	}
	ctx := context.Background()
	services.NewAutomaticUserDataCleanupService().InitAutomaticUserDataCleaner(ctx)

	router := initRouter()
	httpRequestHandler := middlewares.CorsMiddleware(router)
	httpRequestHandler = middlewares.LoggingMiddleware(httpRequestHandler)

	serverPort := configuration.ServerPort.GetValue()
	log.Printf("Starting Server and expose port %d\n", serverPort)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), httpRequestHandler); err != nil {
		log.Fatal("Failed to start Api", err)
	}
}

func initRouter() *mux.Router {
	urlPrefix := "/api/v1.0"
	router := mux.NewRouter()
	router.Path(urlPrefix + "/users/login").HandlerFunc(handlers.PostUsersLogin).Methods(http.MethodPost)
	router.Path(urlPrefix + "/users/me").Handler(withAuth(handlers.GetUsersMe)).Methods(http.MethodGet)
	router.Path(urlPrefix + "/rooms").Handler(withAuth(handlers.GetAllRooms)).Methods(http.MethodGet)
	router.Path(urlPrefix + "/rooms/{id}/bookings").Handler(withAuth(handlers.GetBookingsByRoomAndDate)).Methods(http.MethodGet)

	router.Path(urlPrefix + "/users/{id}/bookings").Handler(withAuth(handlers.GetBookingsByUser)).Methods(http.MethodGet)
	router.Path(urlPrefix + "/bookings").Handler(withAuth(handlers.PostBooking)).Methods(http.MethodPost)
	router.Path(urlPrefix + "/configuration").Handler(withAuth(handlers.GetFrontendConfiguration)).Methods(http.MethodGet)
	return router

}

func withAuth(handlerFunc http.HandlerFunc) http.Handler {
	return middlewares.AuthMiddleware(handlerFunc)
}

func initDefaultConfigurations() {
	configuration.RegisterIntegerConfigurationWithDefault(configuration.ServerPort, 8080)

	configuration.RegisterStringConfigurationWithDefault(configuration.DBHost, "localhost")
	configuration.RegisterStringConfigurationWithDefault(configuration.DBUsername, "postgres")
	configuration.RegisterStringConfigurationWithDefault(configuration.DBPassword, "password123")
	configuration.RegisterStringConfigurationWithDefault(configuration.DBName, "desks_api")

	configuration.RegisterStringConfigurationWithDefault(configuration.JwtSigningKey, "")

	configuration.RegisterIntegerConfigurationWithDefault(configuration.UserdataCleanerIntervalHours, 24)
	configuration.RegisterIntegerConfigurationWithDefault(configuration.MaxUserdataAgeDays, 90)
}
