package middlewares

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println(request.Method, request.URL)
		next.ServeHTTP(writer, request)
	})

}
