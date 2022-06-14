package tracking

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"notifications/internal/api"
	myError "notifications/internal/error"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	var tests1 = []struct {
		name                 string
		url                  string
		code                 string
		repo                 func() APITrackingRepo
		expectedResponseCode int
		expectedBody         api.HttpResponse
	}{
		{
			name: "Should return 200 when tracking code is valid and db return result",
			url:  "http://localhost:8080/api/v1/user/tracking/RandomCode1",
			code: "RandomCode1",
			repo: func() APITrackingRepo {
				response := []Shipment{
					{"0", "3", "testcase01", "testcase01", "testcase01", "testcase01"},
					{"1", "3", "testcase02", "testcase02", "testcase02", "testcase02"},
					{"2", "3", "testcase03", "testcase03", "testcase03", "testcase03"},
					{"3", "3", "testcase04", "testcase04", "testcase04", "testcase04"},
				}
				mockRepo := MockRepo{}
				mockRepo.On("FindAll", "RandomCode1").Return(response, nil).Once()
				return mockRepo
			},
			expectedResponseCode: 200,
			expectedBody: api.HttpResponse{
				Error:   nil,
				Message: "200: All shipment have been found successfully",
				Result: []Shipment{
					{"0", "3", "testcase01", "testcase01", "testcase01", "testcase01"},
					{"1", "3", "testcase02", "testcase02", "testcase02", "testcase02"},
					{"2", "3", "testcase03", "testcase03", "testcase03", "testcase03"},
					{"3", "3", "testcase04", "testcase04", "testcase04", "testcase04"},
				},
			},
		},
		{
			name: "Should return 400 when path parameter is invalid",
			url:  "http://localhost:8080/api/v1/user/tracking/",
			code: "",
			repo: func() APITrackingRepo {
				mockRepo := new(MockRepo)
				return mockRepo
			},
			expectedResponseCode: 400,
			expectedBody: api.HttpResponse{
				Error:   myError.InvalidPara,
				Message: "400: The Parameter is invalid.",
				Result:  nil,
			},
		},
		{
			name: "Should return 500 when db fail",
			url:  "http://localhost:8080/api/v1/user/tracking/RandomCode1",
			code: "RandomCode1",
			repo: func() APITrackingRepo {
				err := myError.NewError(errors.New("db error"), "Database query error.", 500)
				mockRepo := MockRepo{}
				mockRepo.On("FindAll", "RandomCode1").Return([]Shipment{}, err).Once()
				return mockRepo
			},
			expectedResponseCode: 500,
			expectedBody: api.HttpResponse{
				Error:   myError.NewError(errors.New("db error"), "Database query error.", 500),
				Message: "500: Database query error.",
				Result:  nil,
			},
		},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.repo()
			api := NewTrackingHandler(repo)
			r, _ := http.NewRequest("GET", tt.url, nil)
			r = mux.SetURLVars(r, map[string]string{
				"code": tt.code,
			})
			w := httptest.NewRecorder()
			api.FindAll(w, r)

			require.Equal(t, tt.expectedResponseCode, w.Code)
			if &tt.expectedBody != nil {
				expectedJSON, err := json.Marshal(tt.expectedBody)
				require.NoError(t, err)
				require.Equal(t, string(expectedJSON), strings.Trim(string(w.Body.Bytes()), "\n"))
			}
		})
	}
}
