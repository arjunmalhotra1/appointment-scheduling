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
	r.Get("/get-available", h.getAvailableAppointments)
	h.Router = r
	return &h
}

func (h *handler) postAppointment(res http.ResponseWriter, req *http.Request) {
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

func (h *handler) getAvailableAppointments(res http.ResponseWriter, req *http.Request) {
	stringTrainerID := req.URL.Query().Get("trainer_id")
	if stringTrainerID == "" {
		render.Status(req, http.StatusBadRequest)
		render.JSON(res, req, "trainerId missing in the query.")
		return
	}
	stringStartDate := req.URL.Query().Get("starts_at")
	if stringStartDate == "" {
		render.Status(req, http.StatusBadRequest)
		render.JSON(res, req, "starts_at missing in the query.")
		return
	}
	stringEndDate := req.URL.Query().Get("ends_at")
	if stringEndDate == "" {
		render.Status(req, http.StatusBadRequest)
		render.JSON(res, req, "ends_at missing in the query.")
		return
	}

	allSchedules, err := h.scheduler.GetAppointments(stringTrainerID, stringStartDate, stringEndDate)
	if err != nil {
		render.Status(req, http.StatusBadRequest)
		render.JSON(res, req, err.Error())
		return

	}
	if len(allSchedules) == 0 {
		render.JSON(res, req, "no schedules in this time frame")
		return
	}
	render.JSON(res, req, allSchedules)
}
