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

func TestHandler_getActors(t *testing.T) {
	type mockBehavior func(r *mock_service.MockActor)

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
			mockBehavior: func(r *mock_service.MockActor) {
				r.EXPECT().GetActors().Return([]presenter.ActorResponse{
					{Id: 1, Sex: "male", Birthday: "2021-10-12", Name: "username", FilmsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"username\",\"sex\":\"male\",\"birthday\":\"2021-10-12\",\"filmsId\":[1,2]}]\n",
		},
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			mockBehavior: func(r *mock_service.MockActor) {
				r.EXPECT().GetActors().Return([]presenter.ActorResponse{
					{Id: 1, Sex: "male", Birthday: "2021-10-12", Name: "username", FilmsId: []int{1, 2}}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"name\":\"username\",\"sex\":\"male\",\"birthday\":\"2021-10-12\",\"filmsId\":[1,2]}]\n",
		},
		{
			name:                 "Unauthorized",
			mockBehavior:         func(r *mock_service.MockActor) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer USE",
			mockBehavior:         func(r *mock_service.MockActor) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor", pkg.MockJWTAuthUser(handler.actors))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/actor", nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getActor(t *testing.T) {
	type mockBehavior func(r *mock_service.MockActor, id string)

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
			mockBehavior: func(r *mock_service.MockActor, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().GetActor(idd).Return(presenter.ActorResponse{
					Id: 1, Sex: "male", Birthday: "2021-10-12", Name: "username", FilmsId: []int{1, 2}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"username\",\"sex\":\"male\",\"birthday\":\"2021-10-12\",\"filmsId\":[1,2]}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockActor, id string) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			mockBehavior: func(r *mock_service.MockActor, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().GetActor(idd).Return(presenter.ActorResponse{
					Id: 1, Sex: "male", Birthday: "2021-10-12", Name: "username", FilmsId: []int{1, 2}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"username\",\"sex\":\"male\",\"birthday\":\"2021-10-12\",\"filmsId\":[1,2]}\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer USE",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.id)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor/", pkg.MockJWTAuthUser(handler.actor))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/actor/"+test.id, nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_postActor(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockActor, actor presenter.ActorRequest)
	var name = new(string)
	*name = "name"

	var sex = new(string)
	*sex = "sex"

	var birthday = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*birthday = time

	var filmsId = new([]int)
	*filmsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputBody            string
		inputActor           presenter.ActorRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			inputBody:   `{"name": "name", "sex": "sex", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Sex:      sex,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior: func(r *mock_service.MockActor, actor presenter.ActorRequest) {
				r.EXPECT().CreateActor(actor).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "1",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			mockBehavior:         func(r *mock_service.MockActor, actor presenter.ActorRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			mockBehavior:         func(r *mock_service.MockActor, actor presenter.ActorRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.inputActor)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor", pkg.MockJWTAuthAdmin(handler.createActor))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/actor",
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_putActor(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockActor, id string, actor presenter.ActorRequest)
	var name = new(string)
	*name = "name"

	var sex = new(string)
	*sex = "sex"

	var birthday = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*birthday = time

	var filmsId = new([]int)
	*filmsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		inputBody            string
		inputActor           presenter.ActorRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			inputBody:   `{"name": "name", "sex": "sex", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Sex:      sex,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior: func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PutActor(idd, actor).Return(presenter.ActorResponse{
					Id:       1,
					Name:     "name",
					Sex:      "sex",
					Birthday: "2021-10-12",
					FilmsId:  []int{1, 2},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"name\",\"sex\":\"sex\",\"birthday\":\"2021-10-12\",\"filmsId\":[1,2]}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.id, test.inputActor)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor/", pkg.MockJWTAuthAdmin(handler.putActor))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/actor/"+test.id,
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_patchActor(t *testing.T) {
	var dateFormat = "2006-01-02"
	type mockBehavior func(r *mock_service.MockActor, id string, actor presenter.ActorRequest)
	var name = new(string)
	*name = "name"

	var birthday = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*birthday = time

	var filmsId = new([]int)
	*filmsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		inputBody            string
		inputActor           presenter.ActorRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			inputBody:   `{"name": "name", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior: func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PatchActor(idd, actor).Return(presenter.ActorResponse{
					Id:       1,
					Name:     "name",
					Sex:      "sex",
					Birthday: "2021-10-12",
					FilmsId:  []int{1, 2},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"name\",\"sex\":\"sex\",\"birthday\":\"2021-10-12\",\"filmsId\":[1,2]}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:        "Entity not found",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "-1",
			inputBody:   `{"name": "name", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior: func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PatchActor(idd, actor).Return(presenter.ActorResponse{}, errors.New("entity not found"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "entity not found\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string, actor presenter.ActorRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.id, test.inputActor)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor/", pkg.MockJWTAuthAdmin(handler.patchActor))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/actor/"+test.id,
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_deleteActor(t *testing.T) {
	type mockBehavior func(r *mock_service.MockActor, id string)

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
			mockBehavior: func(r *mock_service.MockActor, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().DeleteActor(idd)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockActor, id string) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockActor, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.id)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor/", pkg.MockJWTAuthAdmin(handler.deleteActor))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/actor/"+test.id, nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_actors_invalid_method(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockActor, actor presenter.ActorRequest)
	var name = new(string)
	*name = "name"

	var sex = new(string)
	*sex = "sex"

	var birthday = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*birthday = time

	var filmsId = new([]int)
	*filmsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputBody            string
		inputActor           presenter.ActorRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Method Not Allowed",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			inputBody:   `{"name": "name", "sex": "sex", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Sex:      sex,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior:         func(r *mock_service.MockActor, actor presenter.ActorRequest) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.inputActor)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor", pkg.MockJWTAuthAdmin(handler.actors))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/actor",
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_actor_invalid_method(t *testing.T) {
	var dateFormat = "2006-01-02"

	type mockBehavior func(r *mock_service.MockActor, actor presenter.ActorRequest, id string)
	var name = new(string)
	*name = "name"

	var sex = new(string)
	*sex = "sex"

	var birthday = new(time.Time)
	time, _ := time.Parse(dateFormat, "2021-10-12")
	*birthday = time

	var filmsId = new([]int)
	*filmsId = []int{1, 2}

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		inputBody            string
		id                   string
		method               string
		inputActor           presenter.ActorRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Method Not Allowed",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			method:      "POST",
			id:          "1",
			inputBody:   `{"name": "name", "sex": "sex", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Sex:      sex,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior:         func(r *mock_service.MockActor, actor presenter.ActorRequest, id string) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			mockBehavior: func(r *mock_service.MockActor, actor presenter.ActorRequest, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().GetActor(idd).Return(presenter.ActorResponse{
					Id: 1, Sex: "male", Birthday: "2021-10-12", Name: "username", FilmsId: []int{1, 2}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"name\":\"username\",\"sex\":\"male\",\"birthday\":\"2021-10-12\",\"filmsId\":[1,2]}\n",
		},
		{
			name:        "PUT",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			method:      "PUT",
			inputBody:   `{"name": "name", "sex": "sex", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Sex:      sex,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior:         func(r *mock_service.MockActor, actor presenter.ActorRequest, id string) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:        "PATCH",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			method:      "PATCH",
			inputBody:   `{"name": "name", "birthday": "2021-10-12", "filmsId": [1, 2]}`,
			inputActor: presenter.ActorRequest{
				Name:     name,
				Birthday: birthday,
				FilmsId:  filmsId,
			},
			mockBehavior: func(r *mock_service.MockActor, actor presenter.ActorRequest, id string) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:        "DELETE",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			method: "DELETE",
			mockBehavior: func(r *mock_service.MockActor, actor presenter.ActorRequest, id string) {},
			expectedStatusCode: 403,
			expectedResponseBody: "Forbidden\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockActor(c)
			test.mockBehavior(repo, test.inputActor, test.id)

			services := &service.Service{Actor: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/actor/", pkg.MockJWTAuthAdmin(handler.actor))

			w := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, "/api/actor/"+test.id,
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
