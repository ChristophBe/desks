package util

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func JsonResponseWriter(responseBody interface{}, statusCode int, writer http.ResponseWriter, _ *http.Request) (err error) {

	if responseBody == nil && statusCode == http.StatusNoContent {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	jsonResponse, err := json.Marshal(responseBody)

	if err != nil {
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.WriteHeader(statusCode)
	_, err = writer.Write(jsonResponse)
	return
}

func ErrorResponseWriter(err error, writer http.ResponseWriter, request *http.Request) {
	if handlerError, ok := err.(HandlerError); ok {
		log.Println(handlerError)
		err = JsonResponseWriter(handlerError, handlerError.Status, writer, request)
		if err != nil {
			panic(err)
		}
		return
	}

	log.Println(err)
	err = JsonResponseWriter(InternalServerError(err), http.StatusInternalServerError, writer, request)
	if err != nil {
		panic(err)
	}

}

func GetIntegerUrlParameter(request *http.Request, key string) (_ int64, err error) {
	params := mux.Vars(request)
	value, err := strconv.Atoi(params[key])
	if err != nil {
		return
	}
	return int64(value), err
}
