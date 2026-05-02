package handler

import (
	"L2.18/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// UserHandler обработчик запросов
type UserHandler struct {
	us *service.UserService
}

// ErrResponse структура ответа ошибки
type ErrResponse struct {
	Err string `json:"error"`
}

// ResultResponse структура ответа при успешном выполнении
type ResultResponse struct {
	Result []string `json:"result"`
}

// NewUserHandler создание экземпляра UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{us: userService}
}

// CreateEvent Создание события
func (h *UserHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	h.crudEvent(w, r, "create")
}

// UpdateEvent обновление события
func (h *UserHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	h.crudEvent(w, r, "update")
}

// DeleteEvent удаление события
func (h *UserHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	h.crudEvent(w, r, "delete")
}

// GetEventsDay получение событий на день
func (h *UserHandler) GetEventsDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	h.getEvents(w, r, 1)
}

// GetEventsWeek получение событий на неделю
func (h *UserHandler) GetEventsWeek(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	h.getEvents(w, r, 7)
}

// GetEventsMonth получение событий на месяц вперед
func (h *UserHandler) GetEventsMonth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	h.getEvents(w, r, 30)
}

func (h *UserHandler) crudEvent(w http.ResponseWriter, r *http.Request, operation string) {
	err := r.ParseForm()
	id := r.FormValue("user_id")
	date := r.FormValue("date")
	event := r.FormValue("event")

	switch operation {
	case "create":
		err = h.us.CreateEvent(id, date, event)
	case "delete":
		err = h.us.DeleteEvent(id, date, event)
	case "update":
		err = h.us.UpdateEvent(id, date, event)
	}

	if err != nil {
		sendErrorResponse(w, err)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UserHandler) getEvents(w http.ResponseWriter, r *http.Request, days int) {
	id := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	result, err := h.us.GetEventsForUserID(id, date, days)
	if err != nil {
		sendErrorResponse(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resultResponse := ResultResponse{Result: result}

		err := json.NewEncoder(w).Encode(resultResponse)
		if err != nil {
			fmt.Println(err)
		}

	}
}

func sendErrorResponse(w http.ResponseWriter, err error) {
	statusCode := getHTTPCodeForError(err)
	w.WriteHeader(statusCode)
	errResponse := ErrResponse{Err: err.Error()}

	err = json.NewEncoder(w).Encode(errResponse)
	if err != nil {
		fmt.Println(err)
	}
}

func getHTTPCodeForError(err error) int {

	if errors.Is(err, service.ErrNoSuchEvent) ||
		errors.Is(err, service.ErrInvalidID) ||
		errors.Is(err, service.ErrUserNotFound) {
		return 503
	}

	if errors.Is(err, service.ErrInvalidDate) {
		return 400
	}

	return 500
}
