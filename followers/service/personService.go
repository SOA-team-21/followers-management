package service

import (
	"log"

	"followers.xws.com/model"
	"followers.xws.com/repo"
)

type PersonService struct {
	logger *log.Logger
	repo   *repo.PersonRepo
}

func NewPersonService(l *log.Logger, r *repo.PersonRepo) *PersonService {
	return &PersonService{l, r}
}

func (s *PersonService) GetProfile(userId string) (*model.Person, error) {
	person, err := s.repo.GetPerson(userId)
	if err != nil {
		return nil, err
	}
	return person, nil
}
