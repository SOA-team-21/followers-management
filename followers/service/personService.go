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

func (s *PersonService) Follow(userIdToFollow, userIdFollower string) error {
	err := s.repo.Follow(userIdToFollow, userIdFollower)
	return err
}

func (s *PersonService) UnFollow(userIdToUnFollow, userIdFollower string) error {
	err := s.repo.UnFollow(userIdToUnFollow, userIdFollower)
	return err
}
