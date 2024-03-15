package service

import (
	"filmLibraryVk/internal/repository"
	"filmLibraryVk/model/dto/actor"
)

type ActorService struct {
	repo repository.Actor
}

func NewActorService(repo repository.Actor) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) GetActor(id int) (actor.ActorResponse, error) {
	return s.repo.GetActor(id)
}
func (s *ActorService) GetActors() ([]actor.ActorResponse, error) {
	return s.repo.GetActors()
}

func (s *ActorService) CreateActor(request actor.ActorRequest) (int, error) {
	return s.repo.CreateActor(request)
}

func (s *ActorService) PutActor(id int, request actor.ActorRequest) (actor.ActorResponse, error) {
	return s.repo.PutActor(id, request)
}

func (s *ActorService) PatchActor(id int, request actor.ActorRequest) (actor.ActorResponse, error) {
	return s.repo.PatchActor(id, request)
}

func (s *ActorService) DeleteActor(id int) error {
	return s.repo.DeleteActor(id)
}