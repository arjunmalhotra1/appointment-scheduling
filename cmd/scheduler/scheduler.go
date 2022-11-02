package scheduler

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"
)

type Appointment struct {
	id                         int
	TrainerId                  int64  `json:"trainer_id"`
	UserId                     int64  `json:"user_id"`
	StringAppointmentStartDate string `json:"starts_at"`
	StringAppointmentEndDate   string `json:"ends_at"`
	unixAppointmentStartDate   int64
	unixAppointmentEndDate     int64
}

const timeLayout = "2006-01-02T15:04:05-08:00"

type Scheduler map[int64][]Appointment

func New() *Scheduler {
	s := make(Scheduler)
	return &s
}

func (s Scheduler) AddAppointment(appt Appointment) error {
	apptStartTime, err := time.Parse(timeLayout, appt.StringAppointmentStartDate)
	if err != nil {
		log.Println("AddAppointment: error while parsing startTime")
		return fmt.Errorf("internal server error")
	}
	if int(apptStartTime.Weekday()) > 5 {
		return fmt.Errorf("day should have only be between Monday and Friday")
	}
	if apptStartTime.Hour() < 8 || apptStartTime.Hour() > 17 {
		return fmt.Errorf("first appointment starts at 0800 and last appointment starts at 1630")
	}
	fmt.Println(apptStartTime.Minute())
	if apptStartTime.Minute() != 30 && apptStartTime.Minute() != 0 {
		return fmt.Errorf("start time should have only 0 or 30 minutes")
	}

	apptEndTime, err := time.Parse(timeLayout, appt.StringAppointmentEndDate)
	if err != nil {
		log.Println("AddAppointment: error while parsing endTime")
		return fmt.Errorf("internal server error")
	}
	if int(apptEndTime.Weekday()) > 5 {
		return fmt.Errorf("day should have only be between Monday and Friday")
	}
	if apptEndTime.Hour() < 8 || apptEndTime.Hour() > 17 {
		return fmt.Errorf("first appointment ends at 0830 and last appointment ends at 1700")
	}
	if apptEndTime.Minute() != 30 && apptEndTime.Minute() != 0 {
		return fmt.Errorf("end time should have only 00 or 30 minutes")
	}
	timeDiff := apptEndTime.Sub(apptStartTime)
	if timeDiff.Minutes() != 30 {
		return fmt.Errorf("an appointment should only be 30 minutes of duration")
	}
	appt.id = len(s[appt.TrainerId]) + 1
	appt.unixAppointmentStartDate = apptStartTime.Unix()
	appt.unixAppointmentEndDate = apptEndTime.Unix()

	if _, ok := s[appt.TrainerId]; !ok {
		var tempList = []Appointment{appt}
		s[appt.TrainerId] = tempList
	} else {
		isOverLap := isStartTimeOverlap(s[appt.TrainerId], appt.unixAppointmentStartDate)
		if isOverLap {
			log.Printf("AddAppointment: overlap %v", appt)
			return fmt.Errorf("there was an overlap, this time slot has already been taken")
		}
		s[appt.TrainerId] = append(s[appt.TrainerId], appt)
		sort.Sort(byStartTime(s[appt.TrainerId]))
	}
	s.PrintAppointments()
	return nil
}

type byStartTime []Appointment

func (bST byStartTime) Len() int {
	return len(bST)
}

func (bST byStartTime) Swap(i, j int) {
	bST[i], bST[j] = bST[j], bST[i]
}

func (bST byStartTime) Less(i, j int) bool {
	return bST[i].unixAppointmentStartDate < bST[j].unixAppointmentStartDate
}

func isStartTimeOverlap(appointments []Appointment, startTime int64) bool {
	fmt.Println(appointments, startTime)
	l := 0
	r := len(appointments) - 1
	for l <= r {
		mid := l + (r-l)/2
		if appointments[mid].unixAppointmentStartDate == startTime {
			return true
		} else if appointments[mid].unixAppointmentStartDate < startTime {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return false
}

func (s Scheduler) PrintAppointments() {
	for _, v := range s {
		fmt.Println(v)
	}
}

func (s Scheduler) GetAppointmentsByTrainer(trainerId string) ([]Appointment, error) {
	intTrainerID, err := strconv.ParseInt(trainerId, 10, 64)
	if err != nil {
		return []Appointment{}, fmt.Errorf("GetAppointmentsByTrainer: error converting string to int %v", err)
	}
	return s[intTrainerID], nil
}

func (s Scheduler) GetAppointments(trainerId, stringStartDate, stringEndDate string) ([]Appointment, error) {
	var availableAppointments []Appointment
	startDate, err := time.Parse(timeLayout, stringStartDate)
	if err != nil {
		log.Println("GetAppointments: error parsing the start date", err)
		return []Appointment{}, fmt.Errorf("internal server error")
	}
	endDate, err := time.Parse(timeLayout, stringEndDate)
	if err != nil {
		log.Println("GetAppointments: error parsing the end date", err)
		return []Appointment{}, fmt.Errorf("internal server error")
	}
	if endDate.Before(startDate) {
		return []Appointment{}, fmt.Errorf("start date needs to be before the end date")
	}
	intTrainerId, err := strconv.Atoi(trainerId)
	if err != nil {
		log.Println("GetAppointments: error converting string to int", err)
		return []Appointment{}, fmt.Errorf("internal server error")
	}
	effectiveStartDate := calculateEffectiveStartDate(startDate)
	effectiveEndDate := calculateEffectiveEndDate(endDate)
	fmt.Println("e:", effectiveEndDate)
	for effectiveStartDate.Before(effectiveEndDate) {
		isOverLap := isStartTimeOverlap(s[int64(intTrainerId)], effectiveStartDate.Unix())
		if isOverLap {
			log.Println("Overlap:", effectiveStartDate)
			effectiveStartDate = effectiveStartDate.Add(time.Minute * 30)
			effectiveStartDate = calculateEffectiveStartDate(effectiveStartDate)
			continue
		}
		possibleAppt := Appointment{
			StringAppointmentStartDate: effectiveStartDate.String(),
			StringAppointmentEndDate:   effectiveStartDate.Add(time.Minute * 30).String(),
			TrainerId:                  int64(intTrainerId),
		}
		availableAppointments = append(availableAppointments, possibleAppt)
		effectiveStartDate = effectiveStartDate.Add(time.Minute * 30)
		effectiveStartDate = calculateEffectiveStartDate(effectiveStartDate)

	}
	return availableAppointments, nil
}

func calculateEffectiveStartDate(startDate time.Time) time.Time {
	if startDate.Hour() >= 17 {
		return time.Date(startDate.Year(), startDate.Month(), startDate.Day()+1, 8, 00, 00, 0, startDate.Location())
	} else if startDate.Hour() < 8 {
		return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 8, 00, 00, 0, startDate.Location())
	}
	if startDate.Minute() < 30 && startDate.Minute() > 0 {
		return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour(), 30, 00, 0, startDate.Location())
	} else if startDate.Minute() > 30 {
		return time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour()+1, 00, 00, 0, startDate.Location())
	}
	return startDate

}

func calculateEffectiveEndDate(endDate time.Time) time.Time {
	if endDate.Hour() >= 17 {
		return time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 17, 00, 00, 0, endDate.Location())
	} else if endDate.Hour() < 8 {
		return time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 8, 00, 00, 0, endDate.Location())
	}
	if endDate.Minute() < 30 {
		return time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), 00, 00, 0, endDate.Location())
	}
	return time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), 30, 00, 0, endDate.Location())
}
