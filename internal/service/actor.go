package service

import (
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/internal/repository"
)

type ActorService struct {
	repo repository.Actor
}

func NewActorService(repo repository.Actor) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) GetActor(id int) (presenter.ActorResponse, error) {
	return s.repo.GetActor(id)
}
func (s *ActorService) GetActors() ([]presenter.ActorResponse, error) {
	return s.repo.GetActors()
}

func (s *ActorService) CreateActor(request presenter.ActorRequest) (int, error) {
	return s.repo.CreateActor(request)
}

func (s *ActorService) PutActor(id int, request presenter.ActorRequest) (presenter.ActorResponse, error) {
	return s.repo.PutActor(id, request)
}

func (s *ActorService) PatchActor(id int, request presenter.ActorRequest) (presenter.ActorResponse, error) {
	return s.repo.PatchActor(id, request)
}

func (s *ActorService) DeleteActor(id int) error {
	return s.repo.DeleteActor(id)
}
