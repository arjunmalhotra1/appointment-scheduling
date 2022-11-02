package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/appointment-scheduling/cmd/scheduler"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type handler struct {
	Router    *chi.Mux
	scheduler *scheduler.Scheduler
}

func New() *handler {
	var h handler
	h.scheduler = scheduler.New()
	r := chi.NewMux()
	r.Post("/appointment/", h.postAppointment)
	h.Router = r
	return &h
}

func (h *handler) postAppointment(res http.ResponseWriter, req *http.Request) {
	// TODO: Validate the post body.
	var appt scheduler.Appointment
	err := json.NewDecoder(req.Body).Decode(&appt)
	if err != nil {
		log.Printf("postAppoint: error while decoding json %v", err)
		render.Status(req, http.StatusInternalServerError)
		render.JSON(res, req, "Internal server Error")
		return
	}
	h.scheduler.AddAppointment(appt)
	res.WriteHeader(http.StatusCreated)
	render.JSON(res, req, appt)
}
