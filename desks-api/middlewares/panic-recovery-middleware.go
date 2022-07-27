package middlewares

import (
	"fmt"
	"github.com/ChristophBe/desks/desks-api/util"
	"net/http"
)

func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err, ok := recovered.(error)

				if !ok {
					err = fmt.Errorf("handler panicked; revovered: %v", recovered)
				}

				util.ErrorResponseWriter(util.InternalServerError(err), writer, request)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
