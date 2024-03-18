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
)

func TestHandler_getUsers(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser)

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
			mockBehavior: func(r *mock_service.MockUser) {
				r.EXPECT().GetUsers().Return([]presenter.UserResponse{
					{Id: 1, Username: "username", Role: "USER"}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"username\":\"username\",\"role\":\"USER\"}]\n",
		},
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			mockBehavior: func(r *mock_service.MockUser) {
				r.EXPECT().GetUsers().Return([]presenter.UserResponse{
					{Id: 1, Username: "username", Role: "ADMIN"}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "[{\"id\":1,\"username\":\"username\",\"role\":\"ADMIN\"}]\n",
		},
		{
			name:                 "Unauthorized",
			mockBehavior:         func(r *mock_service.MockUser) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer USE",
			mockBehavior:         func(r *mock_service.MockUser) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/user", pkg.MockJWTAuthUser(handler.mockUsers))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/user", nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, id string)

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
			mockBehavior: func(r *mock_service.MockUser, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().GetUserById(idd).Return(presenter.UserResponse{
					Id: 1, Username: "username", Role: "USER"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"username\":\"username\",\"role\":\"USER\"}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockUser, id string) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			mockBehavior: func(r *mock_service.MockUser, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().GetUserById(idd).Return(presenter.UserResponse{
					Id: 1, Username: "username", Role: "ADMIN"}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"username\":\"username\",\"role\":\"ADMIN\"}\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer USE",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.id)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/user/", pkg.MockJWTAuthUser(handler.mockUser))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/user/"+test.id, nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_putUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, id string, actor presenter.UserRequest)
	var username = new(string)
	*username = "username"

	var password = new(string)
	*password = "password"

	var role = new(string)
	*role = "USER"

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		inputBody            string
		inputUser            presenter.UserRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			inputBody:   `{"username": "username", "password": "password", "role": "USER"}`,
			inputUser: presenter.UserRequest{
				Username: username,
				Password: password,
				Role:     role,
			},
			mockBehavior: func(r *mock_service.MockUser, id string, actor presenter.UserRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PutUser(idd, actor).Return(presenter.UserResponse{
					Id:       1,
					Username: "username",
					Role:     "USER",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"username\":\"username\",\"role\":\"USER\"}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockUser, id string, actor presenter.UserRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string, actor presenter.UserRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string, actor presenter.UserRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.id, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/user/", pkg.MockJWTAuthAdmin(handler.mockUser))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/api/user/"+test.id,
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_patchUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, id string, actor presenter.UserRequest)
	var role = new(string)
	*role = "ADMIN"

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		id                   string
		inputBody            string
		inputUser            presenter.UserRequest
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok admin",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "1",
			inputBody:   `{"role": "ADMIN"}`,
			inputUser: presenter.UserRequest{
				Role: role,
			},
			mockBehavior: func(r *mock_service.MockUser, id string, user presenter.UserRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PatchUser(idd, user).Return(presenter.UserResponse{
					Id:       1,
					Username: "username",
					Role:     "ADMIN",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"username\":\"username\",\"role\":\"ADMIN\"}\n",
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockUser, id string, user presenter.UserRequest) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:        "Entity not found",
			headerName:  "Authorization",
			headerValue: "Bearer ADMIN",
			id:          "-1",
			inputBody:   `{"role": "ADMIN"}`,
			inputUser: presenter.UserRequest{
				Role: role,
			},
			mockBehavior: func(r *mock_service.MockUser, id string, user presenter.UserRequest) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().PatchUser(idd, user).Return(presenter.UserResponse{}, errors.New("entity not found"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "entity not found\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string, user presenter.UserRequest) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string, user presenter.UserRequest) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.id, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/user/", pkg.MockJWTAuthAdmin(handler.mockUser))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PATCH", "/api/user/"+test.id,
				bytes.NewBufferString(test.inputBody))
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_deleteUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, id string)

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
			mockBehavior: func(r *mock_service.MockUser, id string) {
				idd, _ := strconv.Atoi(id)
				r.EXPECT().DeleteUser(idd)
			},
			expectedStatusCode: 200,
		},
		{
			name:                 "Invalid id",
			headerName:           "Authorization",
			headerValue:          "Bearer ADMIN",
			id:                   "1s",
			mockBehavior:         func(r *mock_service.MockUser, id string) {},
			expectedStatusCode:   400,
			expectedResponseBody: "strconv.Atoi: parsing \"1s\": invalid syntax\n",
		},
		{
			name:                 "Forbidden for user",
			headerName:           "Authorization",
			headerValue:          "Bearer USER",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string) {},
			expectedStatusCode:   403,
			expectedResponseBody: "Forbidden\n",
		},
		{
			name:                 "Unauthorized",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockUser, id string) {},
			expectedStatusCode:   401,
			expectedResponseBody: "Invalid JWT token\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.id)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/user/", pkg.MockJWTAuthAdmin(handler.mockUser))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/api/user/"+test.id, nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_users_invalid_method(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Method Not Allowed",
			headerName:  "Authorization",
			headerValue: "Bearer USER",
			mockBehavior: func(r *mock_service.MockUser) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/user", pkg.MockJWTAuthUser(handler.mockUsers))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user", nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_user_invalid_method(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser)

	tests := []struct {
		name                 string
		headerName           string
		headerValue          string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Method Not Allowed",
			headerName:  "Authorization",
			headerValue: "Bearer USER",
			mockBehavior: func(r *mock_service.MockUser) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/user/", pkg.MockJWTAuthUser(handler.mockUser))

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/user/", nil)
			req.Header.Add(test.headerName, test.headerValue)
			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}