package handler

import (
	"errors"
	"github.com/vlad-rubtsov/todo-app-backend/pkg/service"
	mock_service "github.com/vlad-rubtsov/todo-app-backend/pkg/service/mocks"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehaviour func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                string
		headerName          string
		headerValue         string
		token               string
		mockBehaviour       mockBehaviour
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: "1",
		},
		{
			name:        "Empty header",
			headerName:  "",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"empty auth header"}`,
		},
		{
			name:        "Invalid Bearer",
			headerName:  "Authorization",
			headerValue: "Bearr token",
			token:       "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"invalid auth header"}`,
		},
		{
			name:        "Invalid Token",
			headerName:  "Authorization",
			headerValue: "Bearer ",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
			},
			expectedStatusCode:  401,
			expectedRequestBody: `{"message":"token is empty"}`,
		},
		{
			name:        "System Failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehaviour: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, errors.New("failed to parse token"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"message":"failed to parse token"}`,
		},

	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehaviour(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/protected", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}
