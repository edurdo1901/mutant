package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"prueba.com/cmd/api/handlers"
	"prueba.com/internal/mutant"
)

type magnetoMock struct {
	mock.Mock
}

func (m *magnetoMock) IsMutant(ctx context.Context, dna []string) (bool, error) {
	args := m.Called(ctx, dna)
	return args.Get(0).(bool), args.Error(1)
}

func (m *magnetoMock) Stats(ctx context.Context) (mutant.StatsResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(mutant.StatsResponse), args.Error(1)
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

type mockArgs struct {
	methodName string
	inputArgs  []interface{}
	returnArgs []interface{}
	times      int
}

const statsURL = "/stats"
const mutantURL = "/mutant"
const homeURL = "/"

func TestIsMutant(t *testing.T) {

	tt := map[string]struct {
		payload          string
		statusCode       int
		mock             mockArgs
		expectedResponse string
	}{
		"PostMutantTrue": {
			payload:    "test_data/request/mutant_true.json",
			statusCode: http.StatusOK,
			mock: mockArgs{
				methodName: "IsMutant",
				inputArgs:  []interface{}{mock.Anything, mock.Anything},
				returnArgs: []interface{}{true, nil},
			},
		},
		"PostMutantFalse": {
			payload:    "test_data/request/mutant_false.json",
			statusCode: http.StatusForbidden,
			mock: mockArgs{
				methodName: "IsMutant",
				inputArgs:  []interface{}{mock.Anything, mock.Anything},
				returnArgs: []interface{}{false, nil},
			},
		},
		"PostMutantInvalidPayload": {
			payload:          "test_data/request/mutant_invalid_payload.json",
			statusCode:       http.StatusUnprocessableEntity,
			expectedResponse: "test_data/response/mutant_invalid_payload.json",
		},
		"PostMutantErrorProcess": {
			payload:    "test_data/request/mutant_true.json",
			statusCode: http.StatusInternalServerError,
			mock: mockArgs{
				methodName: "IsMutant",
				inputArgs:  []interface{}{mock.Anything, mock.Anything},
				returnArgs: []interface{}{false, errors.New("Error")},
			},
			expectedResponse: "test_data/response/mutant_error_process.json",
		},
		"PostMutantInvalidDNA": {
			payload:    "test_data/request/mutant_invalid_dna.json",
			statusCode: http.StatusUnprocessableEntity,
			mock: mockArgs{
				methodName: "IsMutant",
				inputArgs:  []interface{}{mock.Anything, mock.Anything},
				returnArgs: []interface{}{false, mutant.ErrInvalidDna},
			},
			expectedResponse: "test_data/response/mutant_invalid_dna.json",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			r := SetUpRouter()
			var mMock magnetoMock
			setupMock(tc.mock, &mMock.Mock)
			handler := handlers.New(&mMock)
			handler.API(r)
			req, _ := http.NewRequest("POST", mutantURL, bytes.NewBuffer(readFile(t, tc.payload)))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			responseData, _ := ioutil.ReadAll(w.Body)
			assert.Equal(t, tc.statusCode, w.Code)
			mock.AssertExpectationsForObjects(t, &mMock)
			if tc.expectedResponse != "" {
				assert.JSONEq(t, string(readFile(t, tc.expectedResponse)), string(responseData))
			}

		})
	}
}

func TestStats(t *testing.T) {
	tt := map[string]struct {
		expectedResponse string
		statusCode       int
		mock             mockArgs
	}{
		"GetStats": {
			expectedResponse: "test_data/response/stats_ok.json",
			statusCode:       http.StatusOK,
			mock: mockArgs{
				methodName: "Stats",
				inputArgs:  []interface{}{mock.Anything},
				returnArgs: []interface{}{mutant.StatsResponse{
					CountMutant: 40,
					CountHuman:  100,
					Ratio:       0.4,
				}, nil},
				times: 1,
			},
		},
		"GetStatsError": {
			expectedResponse: "test_data/response/stats_error.json",
			statusCode:       http.StatusInternalServerError,
			mock: mockArgs{
				methodName: "Stats",
				inputArgs:  []interface{}{mock.Anything},
				returnArgs: []interface{}{mutant.StatsResponse{}, errors.New("Error")},
				times:      1,
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			r := SetUpRouter()
			var mMock magnetoMock
			setupMock(tc.mock, &mMock.Mock)
			handler := handlers.New(&mMock)
			handler.API(r)
			req, _ := http.NewRequest("GET", statsURL, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			responseData, _ := ioutil.ReadAll(w.Body)
			assert.JSONEq(t, string(readFile(t, tc.expectedResponse)), string(responseData))
			assert.Equal(t, tc.statusCode, w.Code)
			mock.AssertExpectationsForObjects(t, &mMock)
		})
	}
}

func TestHome(t *testing.T) {
	r := SetUpRouter()
	var mMock magnetoMock
	handler := handlers.New(&mMock)
	handler.API(r)
	req, _ := http.NewRequest("GET", homeURL, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, "\"Prueba mercado libre\"", string(responseData))
	assert.Equal(t, http.StatusOK, w.Code)
}

func readFile(t *testing.T, fileName string) []byte {
	if fileName == "" {
		return []byte("")
	}

	bytes, err := os.ReadFile(fileName)
	require.NoError(t, err)
	return bytes
}

func setupMock(params mockArgs, mock *mock.Mock) {
	if params.methodName != "" {
		mock.On(params.methodName, params.inputArgs...).
			Return(params.returnArgs...).
			Times(1)
	}
}
