package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"followers.xws.com/handler"
	follower "followers.xws.com/proto/followers"
	"followers.xws.com/repo"
	"followers.xws.com/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	lis, err := net.Listen("tcp", ":87")
	if err != nil {
		log.Fatalln(err)
	}

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
	handler := &handler.FollowersHandler{FollowersService: service}

	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			log.Fatalln(err)
		}
	}(lis)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	follower.RegisterFollowersServiceServer(grpcServer, handler)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, syscall.SIGTERM)
	<-stopCh
	grpcServer.Stop()
}

// func main() {
// 	port := os.Getenv("PORT")
// 	if len(port) == 0 {
// 		port = "8080"
// 	}

// 	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	logger := log.New(os.Stdout, "[followers-api]", log.LstdFlags)
// 	repoLogger := log.New(os.Stdout, "[followers-repo]", log.LstdFlags)
// 	serviceLogger := log.New(os.Stdout, "[followers-service]", log.LstdFlags)

// 	repo, err := repo.New(repoLogger)
// 	if err != nil {
// 		logger.Fatal(err)
// 	}
// 	defer repo.CloseDriverConnection(timeoutContext)
// 	repo.CheckConnection()

// 	service := service.NewPersonService(serviceLogger, repo)
// 	handler := handler.NewPersonHandler(service)

// 	router := mux.NewRouter().StrictSlash(true)

// 	router.HandleFunc("/followers/{userId}/profile", handler.GetProfile).Methods("GET")
// 	router.HandleFunc("/followers/{userId}/followers", handler.GetFollowers).Methods("GET")
// 	router.HandleFunc("/followers/{userId}/following", handler.GetFollowing).Methods("GET")
// 	router.HandleFunc("/followers/{userId}/recommended", handler.GetRecommended).Methods("GET")
// 	router.HandleFunc("/followers/{userId}/{followingUserId}/isFollowing", handler.IsFollowing).Methods("GET")
// 	router.HandleFunc("/followers/{toFollow}/{follower}", handler.Follow).Methods("POST")
// 	router.HandleFunc("/followers/{toUnFollow}/{follower}", handler.UnFollow).Methods("DELETE")

// 	cors := gorillaHandlers.CORS(gorillaHandlers.AllowedOrigins([]string{"*"}))

// 	//Initialize the server
// 	server := http.Server{
// 		Addr:    ":" + port,
// 		Handler: cors(router),
// 	}

// 	logger.Println("Server listening on port", port)
// 	//Distribute all the connections to goroutines
// 	go func() {
// 		err := server.ListenAndServe()
// 		if err != nil {
// 			logger.Fatal(err)
// 		}
// 	}()

// 	sigCh := make(chan os.Signal)
// 	signal.Notify(sigCh, os.Interrupt)
// 	signal.Notify(sigCh, os.Kill)

// 	sig := <-sigCh
// 	logger.Println("Received terminate, graceful shutdown", sig)

// 	//Try to shutdown gracefully
// 	if server.Shutdown(timeoutContext) != nil {
// 		logger.Fatal("Cannot gracefully shutdown...")
// 	}
// 	logger.Println("Server stopped")
// }
