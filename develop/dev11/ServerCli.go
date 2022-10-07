package main

import (
	"context"
	"log"
	"net/http"
)

type ServerAPI struct {
	router *http.ServeMux
	server http.Server
}

func InitServerAPI(port string, ec EventController) *ServerAPI {
	router := http.NewServeMux()
	//прописываем маршруты, в обработчик передаем миддлвару
	router.Handle("/create_event", LoggerMiddleware(ec.CreateEvent))
	router.Handle("/update_event", LoggerMiddleware(ec.UpdateEvent))
	router.Handle("/delete_event", LoggerMiddleware(ec.DeleteEvent))
	router.Handle("/events_for_day", LoggerMiddleware(ec.GetEventsForDay))
	router.Handle("/events_for_week", LoggerMiddleware(ec.GetEventsForWeek))
	router.Handle("/events_for_month", LoggerMiddleware(ec.GetEventsForMonth))
	s := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	sApi := &ServerAPI{
		router: router,
		server: s,
	}
	return sApi
}

func (s *ServerAPI) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func (s *ServerAPI) Close() {
	s.server.Shutdown(context.Background())
}
