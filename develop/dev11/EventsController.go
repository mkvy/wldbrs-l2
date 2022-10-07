package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// формат даты
const layout = "2006-01-02"

// реализуем контроллер
type EventController struct {
	service EventService
}

func InitController(s EventService) *EventController {
	return &EventController{
		service: s,
	}
}

// Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month

// CreateEvent POST /create_event
// поля передаются в виде application/x-www-form-urlencoded:
// user_id uint обязательное
// date формата 2006-01-02 обязательное
// message
func (ec *EventController) CreateEvent(w http.ResponseWriter, r *http.Request) {
	//проверяем, пришел ли post method
	logHeader := "createEvent"
	if r.Method != http.MethodPost {
		returnError(w, logHeader, "", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	valUserID := r.FormValue("user_id")
	if valUserID == "" {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseUint(valUserID, 10, 64)
	if err != nil {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	valDate := r.FormValue("date")
	if valDate == "" {
		returnError(w, logHeader, "missing parameter date", http.StatusBadRequest)
		return
	}
	date, err := time.Parse(layout, valDate)
	if err != nil {
		returnError(w, logHeader, "error parsing date", http.StatusBadRequest)
		return
	}
	eventData := r.FormValue("message")
	event := Event{
		ID:        0,
		UserID:    userID,
		Date:      date,
		EventData: eventData,
	}
	err = ec.service.Create(event)
	if err != nil {
		if errors.Is(err, errors.New("event with this ID already exists")) {
			returnError(w, logHeader, err.Error(), http.StatusServiceUnavailable)
		} else {
			returnError(w, logHeader, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	returnResult(w, "event created", http.StatusCreated)
	log.Printf("%s: created event %+v", logHeader, event)
}

// UpdateEvent POST на update_event
// поля передаются в виде application/x-www-form-urlencoded:
// ID uint обязательное
// user_id uint, date, message необязательные
func (ec *EventController) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	//проверяем, пришел ли post method
	logHeader := "updateEvent"
	if r.Method != http.MethodPost {
		returnError(w, logHeader, "", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	//ID - обязательный параметр
	valID := r.FormValue("ID")
	if valID == "" {
		returnError(w, logHeader, "missing/incorrect parameter ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(valID, 10, 64)
	if err != nil {
		returnError(w, logHeader, "missing/incorrect parameter ID", http.StatusBadRequest)
		return
	}

	//ищем event
	event, err := ec.service.Get(id)
	if err != nil {
		if errors.Is(err, errors.New("event not found")) {
			returnError(w, logHeader, err.Error(), http.StatusNotFound)
		} else {
			returnError(w, logHeader, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	//обновляем каждый параметр, если есть

	valUserID := r.FormValue("user_id")

	if valUserID != "" {
		userID, err := strconv.ParseUint(valUserID, 10, 64)
		if err != nil {
			returnError(w, logHeader, "incorrect parameter user_id", http.StatusBadRequest)
			return
		}
		event.UserID = userID
	}
	valDate := r.FormValue("date")

	if valDate != "" {
		date, err := time.Parse(layout, valDate)
		if err != nil {
			returnError(w, logHeader, "error parsing date", http.StatusBadRequest)
			return
		}
		event.Date = date
	}

	eventData := r.FormValue("message")
	if eventData != "" {
		event.EventData = eventData
	}

	err = ec.service.Update(event)
	if err != nil {
		returnError(w, logHeader, err.Error(), http.StatusInternalServerError)
		return
	}
	returnResult(w, "Updated successfully", http.StatusOK)
	log.Printf("%s: updated event %+v", logHeader, event)
}

// DeleteEvent DELETE на delete_event
// поля в виде www-form-urlencoded
// ID - обязательное

func (ec *EventController) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	//проверяем, пришел ли post method
	logHeader := "deleteEvent"
	if r.Method != http.MethodPost {
		returnError(w, logHeader, "", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	//ID - обязательный параметр
	valID := r.FormValue("ID")
	if valID == "" {
		returnError(w, logHeader, "missing/incorrect parameter ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(valID, 10, 64)
	if err != nil {
		returnError(w, logHeader, "missing/incorrect parameter ID", http.StatusBadRequest)
		return
	}
	err = ec.service.Delete(id)
	if err != nil {
		if errors.Is(err, errors.New("event with this ID not found")) {
			returnError(w, logHeader, err.Error(), http.StatusNotFound)
		} else {
			returnError(w, logHeader, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	returnResult(w, "event deleted", http.StatusOK)
	log.Printf("%s: deleted event %+v", logHeader, id)
}

// GET /events_for_day
// поля в форме www-form-urlencoded
// date ,  user_id (тип uint) обязательное

func (ec *EventController) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	logHeader := "getEventsForDay"
	if r.Method != http.MethodGet {
		returnError(w, logHeader, "", http.StatusMethodNotAllowed)
		return
	}
	valDate := r.FormValue("date")
	if valDate == "" {
		returnError(w, logHeader, "missing parameter date", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(layout, valDate)
	if err != nil {
		returnError(w, logHeader, "error parsing date", http.StatusBadRequest)
		return
	}
	valUserID := r.FormValue("user_id")
	if valUserID == "" {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseUint(valUserID, 10, 64)
	if err != nil {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	events := ec.service.GetForDay(userID, date)
	returnEvents(w, logHeader, events)
	log.Printf("%s: get event for user %d", logHeader, userID)
}

//GET events_for_week
// аналогично GET /events_for_day

func (ec *EventController) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	logHeader := "getEventsForWeek"
	if r.Method != http.MethodGet {
		returnError(w, logHeader, "", http.StatusMethodNotAllowed)
		return
	}
	valDate := r.FormValue("date")
	if valDate == "" {
		returnError(w, logHeader, "missing parameter date", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(layout, valDate)
	if err != nil {
		returnError(w, logHeader, "error parsing date", http.StatusBadRequest)
		return
	}
	valUserID := r.FormValue("user_id")
	if valUserID == "" {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseUint(valUserID, 10, 64)
	if err != nil {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	events := ec.service.GetForWeek(userID, date)
	returnEvents(w, logHeader, events)
	log.Printf("%s: get event for user %d", logHeader, userID)
}

// GET events_for_month
// аналогично GET /events_for_day

func (ec *EventController) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	logHeader := "getEventsForMonth"
	if r.Method != http.MethodGet {
		returnError(w, logHeader, "", http.StatusMethodNotAllowed)
		return
	}
	valDate := r.FormValue("date")
	if valDate == "" {
		returnError(w, logHeader, "missing parameter date", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(layout, valDate)
	if err != nil {
		returnError(w, logHeader, "error parsing date", http.StatusBadRequest)
		return
	}
	valUserID := r.FormValue("user_id")
	if valUserID == "" {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.ParseUint(valUserID, 10, 64)
	if err != nil {
		returnError(w, logHeader, "missing/incorrect parameter user_id", http.StatusBadRequest)
		return
	}
	events := ec.service.GetForMonth(userID, date)
	returnEvents(w, logHeader, events)
	log.Printf("%s: get event for user %d", logHeader, userID)
}

// вспомогательные методы writer-ы
// пишут ответ

// пишем результат из одной строки
func returnResult(w http.ResponseWriter, result string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"result": %s}`, result)
}

// возвращаем ошибку
func returnError(w http.ResponseWriter, logHeader, err string, status int) {
	log.Printf("%s: %s", logHeader, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err == "" {
		fmt.Fprintf(w, `{"error": "%s"}`, http.StatusText(status))
		return
	}
	fmt.Fprintf(w, `{"error": "%s: %s"}`, http.StatusText(status), err)
}

// пишем результат из слайса event, сериализуем в json
func returnEvents(w http.ResponseWriter, logHeader string, events []Event) {
	type result struct {
		Result []Event `json:"result"`
	}
	w.Header().Set("Content-Type", "application/json")
	res := result{Result: events}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		returnError(w, logHeader, err.Error(), http.StatusInternalServerError)
		return
	}
}
