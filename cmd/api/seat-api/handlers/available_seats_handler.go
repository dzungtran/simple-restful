package handlers

import (
	"context"
	"log"
	"net/http"
	seat_svc "simple-restful/exmsgs/seat/services"
	"simple-restful/pkg/core/servehttp"
)

type GetAvailableSeatsHandler struct {
	SeatServiceClient seat_svc.SeatServiceClient
}

func (h *GetAvailableSeatsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {

	response, err := h.SeatServiceClient.GetAvailableSeats(context.Background(), &seat_svc.GetAvailableSeatsRequest{})
	if err != nil {
		log.Printf("Error while GetAvailableSeats, details: %v", err.Error())
		servehttp.ResponseJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"message": "Internal server error",
			"status": "fail",
		})
		return
	}

	servehttp.ResponseSuccessJSON(w, response)
	return
}
