package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"followers.xws.com/handler"
	"followers.xws.com/repo"
	"followers.xws.com/service"
	"github.com/gorilla/mux"
)

func startServer(handler *handler.PersonHanlder, port string) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/followers/profile/{userId}", handler.GetProfile).Methods("GET")

	println("Server starting")
	log.Fatal(http.ListenAndServe(port, router))
}

func main() {

	port := ":8080"

	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger := log.New(os.Stdout, "[followers-api]", log.LstdFlags)
	repoLogger := log.New(os.Stdout, "[followers-repo]", log.LstdFlags)
	serviceLogger := log.New(os.Stdout, "[followers-service]", log.LstdFlags)

	repo, err := repo.New(repoLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer repo.CloseDriverConnection(timeoutContext)
	repo.CheckConnection()

	service := service.NewPersonService(serviceLogger, repo)
	handler := handler.NewPersonHandler(service)

	startServer(handler, port) //Port number must be different for different servers (because all run on localhost)
}