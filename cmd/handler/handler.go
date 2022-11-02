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
	r.Post("/appointment", h.postAppointment)
	r.Get("/get-scheduled-appointments", h.getAppointmentsByTrainer)
	//r.Get("/get-available", h.getAppointmentsByAvailability)
	h.Router = r
	return &h
}

func (h *handler) postAppointment(res http.ResponseWriter, req *http.Request) {
	// TODO: Validate the post body.
	// Validate the required fields. Check endTime - startTime = 30mins
	// Validate that both start and end times are always only 00 and 30 mins
	var appt scheduler.Appointment
	err := json.NewDecoder(req.Body).Decode(&appt)
	if err != nil {
		log.Printf("postAppoint: error while decoding json %v", err)
		render.Status(req, http.StatusInternalServerError)
		render.JSON(res, req, "Internal server Error")
		return
	}
	err = h.scheduler.AddAppointment(appt)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		render.JSON(res, req, err.Error())
		return
	}
	res.WriteHeader(http.StatusCreated)
	render.JSON(res, req, appt)
}

func (h *handler) getAppointmentsByTrainer(res http.ResponseWriter, req *http.Request) {
	stringTrainerID := req.URL.Query().Get("trainer_id")
	if stringTrainerID == "" {
		render.Status(req, http.StatusBadRequest)
		render.JSON(res, req, "trainerId missing in the query.")
		return
	}
	allSchedules, err := h.scheduler.GetAppointmentsByTrainer(stringTrainerID)
	if err != nil {
		log.Printf("getAppointmentsByTrainerByDates: error  %v", err)
		render.Status(req, http.StatusInternalServerError)
		render.JSON(res, req, "Internal server Error")
		return
	}
	res.WriteHeader(http.StatusOK)
	render.JSON(res, req, allSchedules)
	//fmt.Println("stringTrainerID: ", stringTrainerID)
}

//func getAppointmentsByAvailability(res http.ResponseWriter, req *http.Request)
