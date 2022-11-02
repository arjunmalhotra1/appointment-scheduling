package scheduler

import "fmt"

type Appointment struct {
	TrainerId int `json:"TrainerId"`
	UserId    int `json:"UserId"`
	// appointmentStartDate time.Time
	// appointmentEndDate   time.Time
}

type Scheduler map[int]Appointment

func New() *Scheduler {
	s := make(Scheduler)
	return &s
}

func (s Scheduler) AddAppointment(appt Appointment) {
	if _, ok := s[appt.TrainerId]; !ok {
		s[appt.TrainerId] = appt
	}
	s.PrintAppointments()
}

func (s Scheduler) PrintAppointments() {
	for _, v := range s {
		fmt.Println(v)
	}
}
