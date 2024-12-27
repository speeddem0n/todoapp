package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	todo "github.com/speeddem0n/todoapp"
	"github.com/speeddem0n/todoapp/pkg/service"
	mock_service "github.com/speeddem0n/todoapp/pkg/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           todo.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"Test","username":"Test","password":"123"}`,
			inputUser: todo.User{
				Name:     "Test",
				Username: "Test",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"username":"Test","password":"123"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Servise Failure",
			inputBody: `{"name":"tt","username":"Test","password":"123"}`,
			inputUser: todo.User{
				Name:     "tt",
				Username: "Test",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			// Test server
			r := gin.New()
			gin.SetMode(gin.ReleaseMode)
			r.POST("/sign-up", handler.singUp)

			// test request

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(testCase.inputBody))

			// Rerform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}

func TestHandler_SingUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, input signInInput)

	testTable := []struct {
		name                string
		inputBody           string
		signInInput         signInInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"username":"Test","password":"123"}`,
			signInInput: signInInput{
				Username: "Test",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("asd123r13tg13fweg352grweaegr", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"asd123r13tg13fweg352grweaegr"}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"password":"123"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, input signInInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"Test","password":"123"}`,
			signInInput: signInInput{
				Username: "Test",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("asd123r13tg13fweg352grweaegr", errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.signInInput)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			// Test server
			r := gin.New()
			gin.SetMode(gin.ReleaseMode)
			r.POST("/sign-in", handler.singIn)

			// test request

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			// Rerform Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}
