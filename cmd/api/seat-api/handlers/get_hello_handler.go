package handlers

import (
	"net/http"
	"simple-restful/pkg/core/servehttp"
)

type GetHelloHandler struct {

}

func (h *GetHelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	servehttp.ResponseSuccessJSON(w, map[string]string{
		"message": "This is a simple restful API",
	})
}
