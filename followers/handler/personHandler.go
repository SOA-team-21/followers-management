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
