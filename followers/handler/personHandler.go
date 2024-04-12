package handler

import (
	"fmt"
	"net/http"

	"followers.xws.com/service"
	"github.com/gorilla/mux"
)

type PersonHanlder struct {
	service *service.PersonService
}

func NewPersonHandler(s *service.PersonService) *PersonHanlder {
	return &PersonHanlder{s}
}

func (p *PersonHanlder) GetProfile(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	userId := vars["userId"]
	if userId == "" {
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}

	person, err := p.service.GetProfile(userId)
	if err != nil {
		return
	}
	if person != nil {
		fmt.Println("This is desired profile: ", person)
	}

	err = person.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (p *PersonHanlder) Follow(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	userIdToFollow := vars["toFollow"]
	userIdFollower := vars["follower"]
	if userIdToFollow == "" {
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}
	if userIdFollower == "" {
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}

	err := p.service.Follow(userIdToFollow, userIdFollower)
	if err != nil {
		http.Error(rw, "Unable to follow", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (p *PersonHanlder) UnFollow(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	userIdToUnFollow := vars["toUnFollow"]
	userIdFollower := vars["follower"]
	if userIdToUnFollow == "" {
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}
	if userIdFollower == "" {
		http.Error(rw, "Unable to convert limit to integer", http.StatusBadRequest)
		return
	}

	err := p.service.Follow(userIdToUnFollow, userIdFollower)
	if err != nil {
		http.Error(rw, "Unable to unfollow", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
