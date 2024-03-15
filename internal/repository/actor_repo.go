package repository

import (
	"database/sql"
	"errors"
	"filmLibraryVk/model/dto/actor"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ActorRepo struct {
	db *sql.DB
}

func NewActorRepo(db *sql.DB) *ActorRepo {
	return &ActorRepo{db: db}
}

func (r *ActorRepo) GetActor(id int) (actor.ActorResponse, error) {
	act := actor.ActorResponse{}
	filmsId := make([]int, 0)
	var birthday string
	var filmId sql.NullInt64

	query, err := r.db.Prepare("SELECT actor.id, name, sex, birthday, film_id FROM actor " +
		"LEFT JOIN actor_film ON actor.id = actor_film.actor_id " +
		"WHERE actor.id = $1")

	if err != nil {
		return actor.ActorResponse{}, err
	}

	defer query.Close()
	rows, err := query.Query(id)

	if err != nil {
		return actor.ActorResponse{}, err
	}

	for rows.Next() {
		err = rows.Scan(&act.Id, &act.Name, &act.Sex, &birthday, &filmId)
		if err != nil {
			return actor.ActorResponse{}, err
		}
		act.Birthday = strings.Split(birthday, "T")[0]
		if filmId.Valid {
			filmsId = append(filmsId, int(filmId.Int64))
		}
	}
	if act.Id != id {
		return actor.ActorResponse{}, errors.New("entity not found")
	}
	act.FilmsId = filmsId
	log.Printf("Get actor with id %d", id)
	return act, nil
}

func (r *ActorRepo) GetActors() ([]actor.ActorResponse, error) {
	mapActor := make(map[int]actor.ActorResponse)
	mapFilms := make(map[int][]int)
	act := actor.ActorResponse{}
	var birthday string
	var filmId sql.NullInt64

	query, err := r.db.Prepare("SELECT actor.id, name, sex, birthday, film_id FROM actor " +
		"LEFT JOIN actor_film ON actor.id = actor_film.actor_id ")
	if err != nil {
		return nil, err
	}
	defer query.Close()
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&act.Id, &act.Name, &act.Sex, &birthday, &filmId)
		if err != nil {
			return nil, err
		}
		act.Birthday = strings.Split(birthday, "T")[0]
		_, ok := mapActor[act.Id]
		if !ok {
			mapActor[act.Id] = act
		}

		_, ok = mapFilms[act.Id]
		if !ok {
			mapFilms[act.Id] = make([]int, 0)
		}
		if filmId.Valid {
			mapFilms[act.Id] = append(mapFilms[act.Id], int(filmId.Int64))
		}
	}
	var actors []actor.ActorResponse
	for key, value := range mapActor {
		value.FilmsId = mapFilms[key]
		actors = append(actors, value)
	}
	log.Printf("Get actors")

	return actors, nil
}

func (r *ActorRepo) CreateActor(request actor.ActorRequest) (int, error) {
	var id int
	query, err := r.db.Prepare("INSERT INTO actor (name, sex, birthday) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer query.Close()
	row, err := query.Query(*request.Name, *request.Sex, *request.Birthday)

	if err != nil {
		return 0, err
	}

	for row.Next() {
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
	}

	query, err = r.db.Prepare("INSERT INTO actor_film (actor_id, film_id) VALUES ($1, $2)")
	if err != nil {
		return 0, err
	}

	for _, val := range *request.FilmsId {
		row, err = query.Query(id, val)

		if err != nil {
			return 0, errors.New("film with such id does not exist")
		}

		row.Next()
	}

	log.Printf("Insert actor with id %d", id)
	return id, nil
}

func (r *ActorRepo) PutActor(id int, request actor.ActorRequest) (actor.ActorResponse, error) {
	var updatedId int
	query, err := r.db.Prepare("UPDATE actor SET name = $1, sex = $2, birthday = $3 WHERE id = $4 RETURNING id")
	if err != nil {
		return actor.ActorResponse{}, err
	}
	defer query.Close()
	row, err := query.Query(*request.Name, *request.Sex, *request.Birthday, id)

	if err != nil {
		return actor.ActorResponse{}, err
	}

	for row.Next() {
		if err := row.Scan(&updatedId); err != nil {
			return actor.ActorResponse{}, err
		}
	}

	if updatedId != id {
		return actor.ActorResponse{}, errors.New("entity not found")
	}

	err = r.updateFilmsId(request, id)

	if err != nil {
		return actor.ActorResponse{}, err
	}

	log.Printf("Put actor with id %d", id)
	return actor.ActorResponse{
		Id:       id,
		Name:     *request.Name,
		Sex:      *request.Sex,
		Birthday: strings.Split((*request.Birthday).String(), " ")[0],
		FilmsId:  *request.FilmsId,
	}, nil
}

func (r *ActorRepo) PatchActor(id int, request actor.ActorRequest) (actor.ActorResponse, error) {
	q := `UPDATE actor SET `
	qParts := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
	var counter = 1
	var updatedId int

	if request.Name != nil {
		qParts = append(qParts, fmt.Sprintf("name=$%d", counter))
		counter++
		args = append(args, request.Name)
	}

	if request.Sex != nil {
		qParts = append(qParts, fmt.Sprintf("sex=$%d", counter))
		counter++
		args = append(args, request.Sex)
	}
	if !request.Birthday.IsZero() {
		qParts = append(qParts, fmt.Sprintf("birthday=$%d", counter))
		counter++
		args = append(args, request.Birthday)
	}
	q += strings.Join(qParts, ",") + ` WHERE id = $` + strconv.Itoa(counter) + "RETURNING id"
	args = append(args, id)

	row, err := r.db.Query(q, args...)

	if err != nil {
		return actor.ActorResponse{}, err
	}

	for row.Next() {
		if err := row.Scan(&updatedId); err != nil {
			return actor.ActorResponse{}, err
		}
	}
	if updatedId != id {
		return actor.ActorResponse{}, errors.New("entity not found")
	}

	err = r.updateFilmsId(request, id)

	if err != nil {
		return actor.ActorResponse{}, err
	}

	log.Printf("Patch actor with id %d", id)
	return r.GetActor(id)
}

func (r *ActorRepo) DeleteActor(id int) error {
	query, err := r.db.Prepare("DELETE FROM actor WHERE id = $1")
	if err != nil {
		return err
	}
	defer query.Close()
	_, err = query.Query(id)

	if err != nil {
		return err
	}
	log.Printf("Delete actor with id %d", id)

	return nil
}

func (r *ActorRepo) updateFilmsId(request actor.ActorRequest, id int) error {
	if request.FilmsId == nil {
		return nil
	}
	query, err := r.db.Prepare("DELETE FROM actor_film WHERE actor_id = $1")

	if err != nil {
		return err
	}
	defer query.Close()

	_, err = query.Query(id)
	if err != nil {
		return err
	}

	query, err = r.db.Prepare("INSERT INTO actor_film (actor_id, film_id) VALUES ($1, $2)")

	if err != nil {
		return err
	}

	for _, val := range *request.FilmsId {
		row, err := query.Query(id, val)

		if err != nil {
			return errors.New("film with such id does not exist")
		}

		row.Next()
	}
	return nil
}