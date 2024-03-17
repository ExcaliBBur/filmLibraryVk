package repository

import (
	"database/sql"
	"errors"
	"filmLibraryVk/api/REST/presenter"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type FilmRepo struct {
	db *sql.DB
}

func NewFilmRepo(db *sql.DB) *FilmRepo {
	return &FilmRepo{db: db}
}

func (r *FilmRepo) GetFilm(id int) (presenter.FilmResponse, error) {
	fil := presenter.FilmResponse{}
	actorsId := make([]int, 0)
	var releaseDate string
	var actorId sql.NullInt64

	query, err := r.db.Prepare("SELECT film.id, name, description, release_date, rating, actor_id FROM film " +
		"LEFT JOIN actor_film ON film.id = actor_film.film_id " +
		"WHERE film.id = $1")

	if err != nil {
		return presenter.FilmResponse{}, err
	}

	defer query.Close()
	row, err := query.Query(id)

	if err != nil {
		return presenter.FilmResponse{}, err
	}

	for row.Next() {
		err = row.Scan(&fil.Id, &fil.Name, &fil.Description, &releaseDate, &fil.Rating, &actorId)
		if err != nil {
			return presenter.FilmResponse{}, err
		}
		fil.ReleaseDate = strings.Split(releaseDate, "T")[0]
		if actorId.Valid {
			actorsId = append(actorsId, int(actorId.Int64))
		}
	}
	if fil.Id != id {
		return presenter.FilmResponse{}, errors.New("entity not found")
	}
	fil.ActorsId = actorsId
	log.Printf("Get film with id %d", id)
	return fil, nil
}

func (r *FilmRepo) GetFilms(sortBy string) ([]presenter.FilmResponse, error) {
	films := make([]presenter.FilmResponse, 0)
	isFilmExistsMap := make(map[int]int)
	mapActors := make(map[int][]int)

	fil := presenter.FilmResponse{}
	var releaseDate string
	var actorId sql.NullInt64

	query, err := r.db.Prepare("SELECT film.id, name, description, release_date, rating, actor_id FROM film " +
		"LEFT JOIN actor_film ON film.id = actor_film.film_id " +
		"ORDER BY " + sortBy)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&fil.Id, &fil.Name, &fil.Description, &releaseDate, &fil.Rating, &actorId)
		if err != nil {
			return nil, err
		}
		fil.ReleaseDate = strings.Split(releaseDate, "T")[0]

		_, ok := mapActors[fil.Id]
		if !ok {
			mapActors[fil.Id] = make([]int, 0)
		}
		if actorId.Valid {
			mapActors[fil.Id] = append(mapActors[fil.Id], int(actorId.Int64))
		}
		_, ok = isFilmExistsMap[fil.Id]
		if !ok {
			films = append(films, fil)
			isFilmExistsMap[fil.Id] = fil.Id
		}
	}

	for i := range films {
		films[i].ActorsId = mapActors[films[i].Id]
	}
	log.Printf("Get films with sort %s", sortBy)
	return films, nil
}

func (r *FilmRepo) CreateFilm(request presenter.FilmRequest) (int, error) {
	var id int
	query, err := r.db.Prepare("INSERT INTO film (name, description, release_date, rating) " +
		"VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer query.Close()
	row, err := query.Query(*request.Name, *request.Description, *request.ReleaseDate, *request.Rating)

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

	for _, val := range *request.ActorsId {
		row, err = query.Query(val, id)

		if err != nil {
			return 0, errors.New("actor with such id does not exist")
		}

		row.Next()
	}

	log.Printf("Insert film with id %d", id)
	return id, nil
}

func (r *FilmRepo) PutFilm(id int, request presenter.FilmRequest) (presenter.FilmResponse, error) {
	var updatedId int
	query, err := r.db.Prepare("UPDATE film SET name = $1, description = $2, release_date = $3, rating = $4" +
		" WHERE id = $5 RETURNING id")
	if err != nil {
		return presenter.FilmResponse{}, err
	}
	defer query.Close()
	row, err := query.Query(*request.Name, *request.Description, *request.ReleaseDate, *request.Rating, id)

	if err != nil {
		return presenter.FilmResponse{}, err
	}

	for row.Next() {
		if err := row.Scan(&updatedId); err != nil {
			return presenter.FilmResponse{}, err
		}
	}
	if updatedId != id {
		return presenter.FilmResponse{}, errors.New("entity not found")
	}

	err = r.updateFilmsId(request, id)

	if err != nil {
		return presenter.FilmResponse{}, err
	}

	log.Printf("Put film with id %d", id)
	return presenter.FilmResponse{
		Id:          id,
		Name:        *request.Name,
		Description: *request.Description,
		ReleaseDate: strings.Split((*request.ReleaseDate).String(), " ")[0],
		Rating:      *request.Rating,
		ActorsId:    *request.ActorsId,
	}, nil
}

func (r *FilmRepo) PatchFilm(id int, request presenter.FilmRequest) (presenter.FilmResponse, error) {
	q := `UPDATE film SET `
	qParts := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
	var counter = 1
	var updatedId int

	if request.Name != nil {
		qParts = append(qParts, fmt.Sprintf("name=$%d", counter))
		counter++
		args = append(args, request.Name)
	}

	if request.Description != nil {
		qParts = append(qParts, fmt.Sprintf("description=$%d", counter))
		counter++
		args = append(args, request.Description)
	}
	if !request.ReleaseDate.IsZero() {
		qParts = append(qParts, fmt.Sprintf("release_date=$%d", counter))
		counter++
		args = append(args, request.ReleaseDate)
	}
	if request.Rating != nil {
		qParts = append(qParts, fmt.Sprintf("rating=$%d", counter))
		counter++
		args = append(args, request.Rating)
	}
	q += strings.Join(qParts, ",") + ` WHERE id = $` + strconv.Itoa(counter) + "RETURNING id"
	args = append(args, id)

	row, err := r.db.Query(q, args...)

	if err != nil {
		return presenter.FilmResponse{}, err
	}

	for row.Next() {
		if err := row.Scan(&updatedId); err != nil {
			return presenter.FilmResponse{}, err
		}
	}
	if updatedId != id {
		return presenter.FilmResponse{}, errors.New("entity not found")
	}

	err = r.updateFilmsId(request, id)

	if err != nil {
		return presenter.FilmResponse{}, err
	}

	log.Printf("Patch film with id %d", id)
	return r.GetFilm(id)
}

func (r *FilmRepo) DeleteFilm(id int) error {
	query, err := r.db.Prepare("DELETE FROM film WHERE id = $1")
	if err != nil {
		return err
	}
	defer query.Close()
	_, err = query.Query(id)

	if err != nil {
		return err
	}
	log.Printf("Delete film with id %d", id)

	return nil
}

func (r *FilmRepo) updateFilmsId(request presenter.FilmRequest, id int) error {
	if request.ActorsId == nil {
		return nil
	}
	query, err := r.db.Prepare("DELETE FROM actor_film WHERE film_id = $1")

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

	for _, val := range *request.ActorsId {
		row, err := query.Query(val, id)

		if err != nil {
			return errors.New("actor with such id does not exist")
		}

		row.Next()
	}
	return nil
}

func (r *FilmRepo) SearchFilmsByName(name string) ([]presenter.FilmResponse, error) {
	films := make([]presenter.FilmResponse, 0)
	isFilmExistsMap := make(map[int]int)
	mapActors := make(map[int][]int)
	fil := presenter.FilmResponse{}
	var releaseDate string
	var actorId sql.NullInt64
	query, err := r.db.Prepare("SELECT film.id, name, description, release_date, rating, actor_id FROM film " +
		"LEFT JOIN actor_film ON film.id = actor_film.film_id " +
		"WHERE film.name LIKE '" + name + "%'")
	if err != nil {
		return nil, err
	}
	defer query.Close()
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&fil.Id, &fil.Name, &fil.Description, &releaseDate, &fil.Rating, &actorId)
		if err != nil {
			return nil, err
		}
		fil.ReleaseDate = strings.Split(releaseDate, "T")[0]

		_, ok := mapActors[fil.Id]
		if !ok {
			mapActors[fil.Id] = make([]int, 0)
		}
		if actorId.Valid {
			mapActors[fil.Id] = append(mapActors[fil.Id], int(actorId.Int64))
		}
		_, ok = isFilmExistsMap[fil.Id]
		if !ok {
			films = append(films, fil)
			isFilmExistsMap[fil.Id] = fil.Id
		}
	}
	for i := range films {
		films[i].ActorsId = mapActors[films[i].Id]
	}
	log.Printf("Search films by name")
	return films, nil
}

func (r *FilmRepo) SearchFilmsByActor(name string) ([]presenter.FilmResponse, error) {
	films := make([]presenter.FilmResponse, 0)
	isFilmExistsMap := make(map[int]int)
	mapActors := make(map[int][]int)

	fil := presenter.FilmResponse{}
	var releaseDate string
	var actorId sql.NullInt64
	query, err := r.db.Prepare("SELECT film.id, film.name, description, release_date, rating, actor_id FROM film " +
		"JOIN actor_film ON film.id = actor_film.film_id " +
		"JOIN actor ON actor_film.actor_id = actor.id " +
		"WHERE actor.name LIKE '" + name + "%'")
	if err != nil {
		return nil, err
	}
	defer query.Close()
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&fil.Id, &fil.Name, &fil.Description, &releaseDate, &fil.Rating, &actorId)
		if err != nil {
			return nil, err
		}
		fil.ReleaseDate = strings.Split(releaseDate, "T")[0]

		_, ok := mapActors[fil.Id]
		if !ok {
			mapActors[fil.Id] = make([]int, 0)
		}
		if actorId.Valid {
			mapActors[fil.Id] = append(mapActors[fil.Id], int(actorId.Int64))
		}
		_, ok = isFilmExistsMap[fil.Id]
		if !ok {
			films = append(films, fil)
			isFilmExistsMap[fil.Id] = fil.Id
		}
	}
	for i := range films {
		films[i].ActorsId = mapActors[films[i].Id]
	}
	log.Printf("Search film by actor")
	return films, nil
}
