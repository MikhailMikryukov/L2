package main

import (
	"L2.18/handler"
	"L2.18/middleware"
	"errors"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "port")
	flag.Parse()

	router := handler.NewRouter()
	loggedRouter := middleware.Logging(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: loggedRouter,
	}

	// Запускаем сервер на порту
	fmt.Println("Starting server at port", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Could not start server: %v\n", err)
	}

}
