package handler

import (
	"L2.18/service"
	"net/http"
)

// NewRouter создание и настройка обработчика
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	userService := service.NewUserService()
	userHandler := NewUserHandler(userService)

	mux.HandleFunc("/create_event", userHandler.CreateEvent)
	mux.HandleFunc("/update_event", userHandler.UpdateEvent)
	mux.HandleFunc("/delete_event", userHandler.DeleteEvent)
	mux.HandleFunc("/events_for_day", userHandler.GetEventsDay)
	mux.HandleFunc("/events_for_week", userHandler.GetEventsWeek)
	mux.HandleFunc("/events_for_month", userHandler.GetEventsMonth)

	return mux
}
