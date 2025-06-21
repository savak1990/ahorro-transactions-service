package handler

import "net/http"

type CommonHandler interface {
	HandleHealth(http.ResponseWriter, *http.Request)
	HandleInfo(http.ResponseWriter, *http.Request)
}
