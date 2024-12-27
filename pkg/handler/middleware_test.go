package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/speeddem0n/todoapp/pkg/service"
	mock_service "github.com/speeddem0n/todoapp/pkg/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct { // Тблица с тестовыми данными
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "1",
		},
		{
			name:                "Invalid Header Name",
			headerName:          "",
			headerValue:         "Bearer token",
			token:               "token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"empty auth header"}`,
		},
		{
			name:                "Invalid Header Value",
			headerName:          "Authorization",
			headerValue:         "Bear token",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"invalid auth header"}`,
		},
		{
			name:                "Invalid Token",
			headerName:          "Authorization",
			headerValue:         "Bearer ",
			mockBehavior:        func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"token is empty"}`,
		},
		{
			name:        "invalid token",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, errors.New("invalid token"))
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"invalid token"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			service := &service.Service{Authorization: auth}
			handler := NewHandler(service)

			// test server

			r := gin.New()
			r.GET("/identity", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			// Make request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestGetUserId(t *testing.T) {
	var getContext = func(id int) *gin.Context {
		ctx := &gin.Context{}
		ctx.Set(userCtx, id)
		return ctx
	}

	testTable := []struct { // Тблица с тестовыми данными
		name       string
		ctx        *gin.Context
		id         int
		shouldFail bool
		err        error
	}{
		{
			name: "OK",
			ctx:  getContext(1),
			id:   1,
		}, {

			name:       "Empty",
			ctx:        &gin.Context{},
			shouldFail: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			id, err := getUserId(testCase.ctx)
			if err != nil {
				if testCase.shouldFail {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
			assert.Equal(t, testCase.id, id)
		})

	}
}
