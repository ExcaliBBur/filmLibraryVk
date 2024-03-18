package handler

import (
	"bytes"
	"errors"
	"filmLibraryVk/api/REST/presenter"
	"filmLibraryVk/internal/service"
	mock_service "filmLibraryVk/internal/service/mocks"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_register(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, user presenter.Register)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            presenter.Register
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "password"}`,
			inputUser: presenter.Register{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Register) {
				r.EXPECT().Register(user).Return("test_token", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: "jwt: test_token",
		},
		{
			name:      "Min length username fail register",
			inputBody: `{"username": "1", "password": "password"}`,
			inputUser: presenter.Register{
				Username: "1",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Register) {},
			expectedStatusCode:   400,
			expectedResponseBody: "Invalid request body. Username length must be >= 2, password length must be [8, 16]\n",
		},
		{
			name:      "Min length password fail register",
			inputBody: `{"username": "username", "password": "passwor"}`,
			inputUser: presenter.Register{
				Username: "1",
				Password: "passwor",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Register) {},

			expectedStatusCode:   400,
			expectedResponseBody: "Invalid request body. Username length must be >= 2, password length must be [8, 16]\n",
		},
		{
			name:      "Max length password fail register",
			inputBody: `{"username": "username", "password": "passwordpasswordpassword"}`,
			inputUser: presenter.Register{
				Username: "1",
				Password: "passwordpasswordpassword",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Register) {},
			expectedStatusCode:   400,
			expectedResponseBody: "Invalid request body. Username length must be >= 2, password length must be [8, 16]\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/auth/register", http.HandlerFunc(handler.register))

			w := httptest.NewRecorder()
			req:= httptest.NewRequest("POST", "/api/auth/register",
				bytes.NewBufferString(test.inputBody))

			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_register_invalid_method(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, user presenter.Register)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            presenter.Register
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Method Not Allowed",
			inputBody: `{"username": "username", "password": "password"}`,
			inputUser: presenter.Register{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Register) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/auth/register", http.HandlerFunc(handler.register))

			w := httptest.NewRecorder()
			req:= httptest.NewRequest("PUT", "/api/auth/register",
				bytes.NewBufferString(test.inputBody))

			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_authenticate(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, user presenter.Login)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            presenter.Login
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"username": "username", "password": "password"}`,
			inputUser: presenter.Login{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Login) {
				r.EXPECT().Login(user).Return("test_token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "jwt: test_token",
		},
		{
			name:      "Invalid credentials",
			inputBody: `{"username": "usernam", "password": "password"}`,
			inputUser: presenter.Login{
				Username: "usernam",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Login) {
				r.EXPECT().Login(user).Return("", errors.New("Invalid login or password"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "Invalid login or password\n",
		},
		{
			name:      "Min length username fail login",
			inputBody: `{"username": "1", "password": "password"}`,
			inputUser: presenter.Login{
				Username: "1",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Login) {},
			expectedStatusCode:   400,
			expectedResponseBody: "Invalid request body. Username length must be >= 2, password length must be [8, 16]\n",
		},
		{
			name:      "Min length password fail login",
			inputBody: `{"username": "username", "password": "passwor"}`,
			inputUser: presenter.Login{
				Username: "1",
				Password: "passwor",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Login) {},

			expectedStatusCode:   400,
			expectedResponseBody: "Invalid request body. Username length must be >= 2, password length must be [8, 16]\n",
		},
		{
			name:      "Max length password fail login",
			inputBody: `{"username": "username", "password": "passwordpasswordpassword"}`,
			inputUser: presenter.Login{
				Username: "1",
				Password: "passwordpasswordpassword",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Login) {},
			expectedStatusCode:   400,
			expectedResponseBody: "Invalid request body. Username length must be >= 2, password length must be [8, 16]\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/auth/authenticate", http.HandlerFunc(handler.authenticate))

			w := httptest.NewRecorder()
			req:= httptest.NewRequest("POST", "/api/auth/authenticate",
				bytes.NewBufferString(test.inputBody))

			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_authenticate_invalid_method(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, user presenter.Login)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            presenter.Login
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Method Not Allowed",
			inputBody: `{"username": "username", "password": "password"}`,
			inputUser: presenter.Login{
				Username: "username",
				Password: "password",
			},
			mockBehavior: func(r *mock_service.MockUser, user presenter.Login) {},
			expectedStatusCode:   405,
			expectedResponseBody: "Method Not Allowed\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{User: repo}
			handler := Handler{services}

			mux := http.NewServeMux()

			mux.Handle("/api/auth/authenticate", http.HandlerFunc(handler.authenticate))

			w := httptest.NewRecorder()
			req:= httptest.NewRequest("PUT", "/api/auth/authenticate",
				bytes.NewBufferString(test.inputBody))

			mux.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}