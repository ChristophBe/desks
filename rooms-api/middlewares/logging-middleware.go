package middlewares

import (
	"log"
	"net/http"
)

func LoggingMiddleware(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.Method, request.URL)
		handler.ServeHTTP(writer, request)
	})

}
