package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/speeddem0n/todoapp/internal/models"
	"github.com/speeddem0n/todoapp/internal/service"
	mock_service "github.com/speeddem0n/todoapp/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user models.User)

	testTable := []struct { // Тблица с тестовыми данными
		name                string
		inputBody           string
		inputUser           models.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"Test","username":"Test","password":"123"}`,
			inputUser: models.User{
				Name:     "Test",
				Username: "Test",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Wrong Input",
			inputBody:           `{"username":"Test","password":"123"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user models.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Servise Error",
			inputBody: `{"name":"tt","username":"Test","password":"123"}`,
			inputUser: models.User{
				Name:     "tt",
				Username: "Test",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user models.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("something went wrong"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"something went wrong"}`,
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
			r.POST("/sign-up", handler.signUp)

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

	testTable := []struct { // Тблица с тестовыми данными
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
			name:                "Wrong Input",
			inputBody:           `{"password":"123"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, input signInInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Servise Error",
			inputBody: `{"username":"Test","password":"123"}`,
			signInInput: signInInput{
				Username: "Test",
				Password: "123",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input signInInput) {
				s.EXPECT().GenerateToken(input.Username, input.Password).Return("asd123r13tg13fweg352grweaegr", errors.New("something went wrong"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"something went wrong"}`,
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
			r.POST("/sign-in", handler.signIn)

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
