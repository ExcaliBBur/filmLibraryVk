package repository

import (
	"database/sql"
	"errors"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/internal/model/entity"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetUserById(id int) (presenter.UserResponse, error) {
	_user := presenter.UserResponse{}

	query, err := r.db.Prepare("SELECT _user.id, _user.username, role FROM _user " +
		"JOIN role ON _user.role_id = role.id " +
		"WHERE _user.id = $1")

	if err != nil {
		return presenter.UserResponse{}, err
	}

	defer query.Close()
	row, err := query.Query(id)

	if err != nil {
		return presenter.UserResponse{}, err
	}

	for row.Next() {
		err = row.Scan(&_user.Id, &_user.Username, &_user.Role)
		if err != nil {
			return presenter.UserResponse{}, err
		}
	}
	if _user.Id != id {
		return presenter.UserResponse{}, errors.New("entity not found")
	}
	log.Printf("Get user with id %d", id)
	return _user, nil
}

func (r *UserRepo) GetUserByUsername(username string) (entity.User, error) {
	_user := entity.User{}

	query, err := r.db.Prepare("SELECT * FROM _user " +
		"WHERE _user.username = $1")

	if err != nil {
		return entity.User{}, err
	}

	defer query.Close()
	row, err := query.Query(username)

	if err != nil {
		return entity.User{}, err
	}

	for row.Next() {
		err = row.Scan(&_user.Id, &_user.Username, &_user.Password, &_user.RoleId)
		if err != nil {
			return entity.User{}, err
		}
	}
	if _user.Username != username {
		return entity.User{}, errors.New("entity not found")
	}
	log.Printf("Get user with username %s", username)
	return _user, nil
}

func (r *UserRepo) GetUsers() ([]presenter.UserResponse, error) {
	users := make([]presenter.UserResponse, 0)
	_user := presenter.UserResponse{}

	query, err := r.db.Prepare("SELECT _user.id, _user.username, role FROM _user " +
		"JOIN role ON _user.role_id = role.id ")

	if err != nil {
		return nil, err
	}

	defer query.Close()
	row, err := query.Query()

	if err != nil {
		return nil, err
	}

	for row.Next() {
		err = row.Scan(&_user.Id, &_user.Username, &_user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, _user)
	}
	log.Printf("Get users")
	return users, nil
}

func (r *UserRepo) CreateUser(register presenter.Register) (int, error) {
	var id int
	query, err := r.db.Prepare(`INSERT INTO _user (username, password, role_id) VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return 0, err
	}
	defer query.Close()
	row, err := query.Query(register.Username, register.Password, 1)

	if err != nil {
		return 0, err
	}

	for row.Next() {
		if err := row.Scan(&id); err != nil {
			return 0, err
		}
	}

	log.Printf("Register user with id %d", id)
	return id, nil
}

func (r *UserRepo) PutUser(id int, request presenter.UserRequest) (presenter.UserResponse, error) {
	var updatedId int
	query, err := r.db.Prepare("UPDATE _user SET username = $1, password = $2, role_id = $3 WHERE id = $4 RETURNING id")
	if err != nil {
		return presenter.UserResponse{}, err
	}
	defer query.Close()
	var role int

	switch *request.Role {
	case "ADMIN":
		role = 1
	case "USER":
		role = 2
	default:
		return presenter.UserResponse{}, errors.New("Role " + *request.Role +
			" does not exist. Available roles: USER, ADMIN")
	}

	row, err := query.Query(*request.Username, *request.Password, role, id)

	if err != nil {
		return presenter.UserResponse{}, err
	}

	for row.Next() {
		if err := row.Scan(&updatedId); err != nil {
			return presenter.UserResponse{}, err
		}
	}

	if updatedId != id {
		return presenter.UserResponse{}, errors.New("entity not found")
	}

	log.Printf("Put user with id %d", id)
	return presenter.UserResponse{
		Id:       id,
		Username: *request.Username,
		Role:     *request.Role,
	}, nil
}

func (r *UserRepo) PatchUser(id int, request presenter.UserRequest) (presenter.UserResponse, error) {
	q := `UPDATE _user SET `
	qParts := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
	var counter = 1
	var updatedId int

	if request.Username != nil {
		qParts = append(qParts, fmt.Sprintf("username=$%d", counter))
		counter++
		args = append(args, request.Username)
	}

	if request.Password != nil {
		qParts = append(qParts, fmt.Sprintf("password=$%d", counter))
		counter++
		args = append(args, request.Password)
	}
	if request.Role != nil {
		qParts = append(qParts, fmt.Sprintf("role_id=$%d", counter))
		counter++
		var role int

		switch *request.Role {
		case "ADMIN":
			role = 1
		case "USER":
			role = 2
		default:
			return presenter.UserResponse{}, errors.New("Role " + *request.Role +
				" does not exist. Available roles: USER, ADMIN")
		}
		args = append(args, role)
	}
	q += strings.Join(qParts, ",") + ` WHERE id = $` + strconv.Itoa(counter) + "RETURNING id"
	args = append(args, id)

	row, err := r.db.Query(q, args...)

	if err != nil {
		return presenter.UserResponse{}, errors.New("User with such username already exists")
	}

	for row.Next() {
		if err := row.Scan(&updatedId); err != nil {
			return presenter.UserResponse{}, err
		}
	}
	if updatedId != id {
		return presenter.UserResponse{}, errors.New("entity not found")
	}

	log.Printf("Patch user with id %d", id)
	return r.GetUserById(id)
}

func (r *UserRepo) DeleteUser(id int) error {
	query, err := r.db.Prepare("DELETE FROM _user WHERE id = $1")
	if err != nil {
		return err
	}
	defer query.Close()
	_, err = query.Query(id)

	if err != nil {
		return err
	}
	log.Printf("Delete user with id %d", id)

	return nil
}
