package handler

import (
	"bytes"
	"errors"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/internal/service"
	mock_service "filmLibraryVk/internal/service/mocks"
	"filmLibraryVk/pkg"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestHandler_getFilms(t *testing.T) {
	type mockBehavior func(r *mock_service.MockFilm)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok user",
			headerName:  "Authorization",
			headerValue: "Bearer USER",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().GetFilms("").Return([]presenter.FilmResponse{
					{Id: 1, Name: "name", Description: "description",
						ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"name\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}]\n",
		},
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().GetFilms("").Return([]presenter.FilmResponse{
					{Id: 1, Name: "name", Description: "description",
						ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"name\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}]\n",
		},
		{
			name:                 "Unauthorized",
			mockBehavior:         func(r *mock_service.MockFilm) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer USE",
			mockBehavior:         func(r *mock_service.MockFilm) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film", pkg.MockJWTAuthUser(handler.mockFilms))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/film", nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getFilm(t *testing.T) {
	type mockBehavior func(r *mock_service.MockFilm, id string)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok user",
			headerName:  "Authorization",
			headerValue: "Bearer USER",
			id:          "1",
			mockBehavior: func(r *mock_service.MockFilm, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().GetFilm(idd).Return(presenter.FilmResponse{
					Id: 1, Name: "name", Description: "description",
					ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"name\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockFilm, id string) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			mockBehavior: func(r *mock_service.MockFilm, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().GetFilm(idd).Return(presenter.FilmResponse{
					Id: 1, Name: "name", Description: "description",
					ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"name\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer USE",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.id)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film/", pkg.MockJWTAuthUser(handler.mockFilm))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/film/"+test.id, nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_postFilm(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockFilm, actor presenter.FilmRequest)
	var name = new(string)
	*name = "name"

	var description = new(string)
	*description = "description"

	var releaseDate = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*releaseDate = time

	var rating = new(int)
	*rating = 5

	var actorsId = new([]int)
	*actorsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputBody            string
		inputFilm            presenter.FilmRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			inputBody: `{"name": "name", "description": "description", 
						"releaseDate": "2021-10-12", "rating": 5, "actorsId": [1, 2]}`,
			inputFilm: presenter.FilmRequest{
				Name:        name,
				Description: description,
				ReleaseDate: releaseDate,
				Rating:      rating,
				ActorsId:    actorsId,
			},
			mockBehavior: func(r *mock_service.MockFilm, film presenter.FilmRequest) {
				r.EXPECT().CreateFilm(film).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "1",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			mockBehavior:         func(r *mock_service.MockFilm, film presenter.FilmRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			mockBehavior:         func(r *mock_service.MockFilm, film presenter.FilmRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.inputFilm)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film", pkg.MockJWTAuthAdmin(handler.mockFilms))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/film",
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_putFilm(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockFilm, id string, actor presenter.FilmRequest)
	var name = new(string)
	*name = "name"

	var description = new(string)
	*description = "description"

	var releaseDate = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*releaseDate = time

	var rating = new(int)
	*rating = 5

	var actorsId = new([]int)
	*actorsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		inputBody            string
		inputFilm            presenter.FilmRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			inputBody: `{"name": "name", "description": "description", 
						"releaseDate": "2021-10-12", "rating": 5, "actorsId": [1, 2]}`,
			inputFilm: presenter.FilmRequest{
				Name:        name,
				Description: description,
				ReleaseDate: releaseDate,
				Rating:      rating,
				ActorsId:    actorsId,
			},
			mockBehavior: func(r *mock_service.MockFilm, id string, actor presenter.FilmRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PutFilm(idd, actor).Return(presenter.FilmResponse{
					Id:          1,
					Name:        "name",
					Description: "description",
					ReleaseDate: "2021-10-12",
					Rating:      5,
					ActorsId:    []int{1, 2},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"name\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockFilm, id string, actor presenter.FilmRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string, actor presenter.FilmRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string, actor presenter.FilmRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.id, test.inputFilm)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film/", pkg.MockJWTAuthAdmin(handler.mockFilm))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/film/"+test.id,
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_patchFilm(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockFilm, id string, actor presenter.FilmRequest)
	var name = new(string)
	*name = "name"

	var rating = new(int)
	*rating = 5

	var releaseDate = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*releaseDate = time

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		inputBody            string
		inputFilm            presenter.FilmRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			inputBody:   `{"name": "name", "rating": 5, "releaseDate": "2021-10-12"}`,
			inputFilm: presenter.FilmRequest{
				Name:        name,
				Rating:      rating,
				ReleaseDate: releaseDate,
			},
			mockBehavior: func(r *mock_service.MockFilm, id string, film presenter.FilmRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PatchFilm(idd, film).Return(presenter.FilmResponse{
					Name:        "name",
					Description: "description",
					ReleaseDate: "2021-10-12",
					Rating:      5,
					ActorsId:    []int{1, 2},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":0,\"name\":\"name\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockFilm, id string, film presenter.FilmRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:        "Entity not found",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "-1",
			inputBody:   `{"name": "name", "rating": 5, "releaseDate": "2021-10-12"}`,
			inputFilm: presenter.FilmRequest{
				Name:        name,
				Rating:      rating,
				ReleaseDate: releaseDate,
			},
			mockBehavior: func(r *mock_service.MockFilm, id string, film presenter.FilmRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PatchFilm(idd, film).Return(presenter.FilmResponse{}, errors.New("entity not found"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "entity not found\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string, film presenter.FilmRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string, film presenter.FilmRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.id, test.inputFilm)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film/", pkg.MockJWTAuthAdmin(handler.mockFilm))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/film/"+test.id,
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_deleteFilm(t *testing.T) {
	type mockBehavior func(r *mock_service.MockFilm, id string)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			mockBehavior: func(r *mock_service.MockFilm, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().DeleteFilm(idd)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockFilm, id string) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockFilm, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.id)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film/", pkg.MockJWTAuthAdmin(handler.mockFilm))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/film/"+test.id, nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_films_invalid_method(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockFilm, actor presenter.FilmRequest)
	var name = new(string)
	*name = "name"

	var description = new(string)
	*description = "description"

	var releaseDate = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*releaseDate = time

	var rating = new(int)
	*rating = 5

	var actorsId = new([]int)
	*actorsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputBody            string
		inputFilm            presenter.FilmRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Method Not Allowed",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			inputBody: `{"name": "name", "description": "description", 
						"releaseDate": "2021-10-12", "rating": 5, "actorsId": [1, 2]}`,
			inputFilm: presenter.FilmRequest{
				Name:        name,
				Description: description,
				ReleaseDate: releaseDate,
				Rating:      rating,
				ActorsId:    actorsId,
			},
			mockBehavior:         func(r *mock_service.MockFilm, film presenter.FilmRequest) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.inputFilm)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film", pkg.MockJWTAuthAdmin(handler.mockFilms))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/film",
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_film_invalid_method(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockFilm, actor presenter.FilmRequest)
	var name = new(string)
	*name = "name"

	var description = new(string)
	*description = "description"

	var releaseDate = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*releaseDate = time

	var rating = new(int)
	*rating = 5

	var actorsId = new([]int)
	*actorsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputBody            string
		inputFilm            presenter.FilmRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Method Not Allowed",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			inputBody: `{"name": "name", "description": "description", 
						"releaseDate": "2021-10-12", "rating": 5, "actorsId": [1, 2]}`,
			inputFilm: presenter.FilmRequest{
				Name:        name,
				Description: description,
				ReleaseDate: releaseDate,
				Rating:      rating,
				ActorsId:    actorsId,
			},
			mockBehavior:         func(r *mock_service.MockFilm, film presenter.FilmRequest) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo, test.inputFilm)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film/", pkg.MockJWTAuthAdmin(handler.mockFilm))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/film/",
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_searchFilms(t *testing.T) {
	type mockBehavior func(r *mock_service.MockFilm)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		field                string
		value                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Search by name for user",
			headerName:  "Authorization",
			headerValue: "Bearer USER",
			field:       "name",
			value:       "1",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().SearchFilmsBy("name", "1").Return([]presenter.FilmResponse{
					{Id: 1, Name: "1", Description: "description",
						ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"1\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}]\n",
		},
		{
			name:        "Search by name for admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			field:       "name",
			value:       "1",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().SearchFilmsBy("name", "1").Return([]presenter.FilmResponse{
					{Id: 1, Name: "1", Description: "description",
						ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"1\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}]\n",
		},
		{
			name:        "Search by actor for user",
			headerName:  "Authorization",
			headerValue: "Bearer USER",
			field:       "actor",
			value:       "actorName",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().SearchFilmsBy("actor", "actorName").Return([]presenter.FilmResponse{
					{Id: 1, Name: "1", Description: "description",
						ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"1\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}]\n",
		},
		{
			name:        "Search by actor for admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			field:       "actor",
			value:       "actorName",
			mockBehavior: func(r *mock_service.MockFilm) {
				r.EXPECT().SearchFilmsBy("actor", "actorName").Return([]presenter.FilmResponse{
					{Id: 1, Name: "1", Description: "description",
						ReleaseDate: "2021-10-12", Rating: 5, ActorsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"1\",\"description\":\"description\",\"ReleaseDate\":\"2021-10-12\",\"Rating\":5,\"actorsId\":[1,2]}]\n",
		},
		{
			name:                 "Unauthorized",
			mockBehavior:         func(r *mock_service.MockFilm) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer USE",
			mockBehavior:         func(r *mock_service.MockFilm) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockFilm(c)
			test.mockBehavior(repo)

			services := &service.Service{Film: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/film/search", pkg.MockJWTAuthUser(handler.searchFilms))
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/film/search", nil)
			req.Header.Add(test.headerName, test.headerValue)
			q := req.URL.Query()
			q.Add(test.field, test.value)
			req.URL.RawQuery = q.Encode()

			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
