package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yervsil/auth_service/domain"
	"github.com/yervsil/auth_service/internal/transport/http/mocks"
	"github.com/yervsil/auth_service/internal/utils"
)


func TestSign_up(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	h := &Handler{
		service: mockService,
		log:     log,
	}

	tests := []struct {
		name               string
		inputBody          any
		mockEncryptPassword func()
		expectedStatusCode int
		expectedResponse   utils.Response
	}{
		{
			name: "valid signup request",
			inputBody: domain.SignupRequest{
				Name: 		"kuku",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockEncryptPassword: func() {
				mockService.EXPECT().EncryptPassword(gomock.Any()).Return(1, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   utils.Response{Message:1},
		},
		{
			name: "invalid request body (validation error)",
			inputBody: domain.SignupRequest{
				Email:    "",
				Password: "password123",
			},
			mockEncryptPassword: func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   utils.Response{Message:[]interface {}{"Field 'Name' failed validation with tag 'required'", "Field 'Email' failed validation with tag 'required'"}},
		},
		{
			name: "error reading request body",
			inputBody: nil,
			mockEncryptPassword: func() {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   utils.Response{Message: "unexpected end of JSON input"},
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var body []byte
			if tt.inputBody != nil {
				body, _ = json.Marshal(tt.inputBody)
			}else{
				body = nil
			}
		
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
			w := httptest.NewRecorder()

			tt.mockEncryptPassword()

			// Act
			h.Sign_up()(w, req)

			// Assert
			res := w.Result()
			defer res.Body.Close()

			responseBody, _ := io.ReadAll(res.Body)
			
			var actualResponse utils.Response
			json.Unmarshal(responseBody, &actualResponse)

			if _, ok := actualResponse.Message.(float64); ok {
				actualResponse.Message = int(actualResponse.Message.(float64))
			}

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tt.expectedResponse, actualResponse)
		})
	}
}

