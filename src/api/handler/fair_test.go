package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/marcovargas74/m74-fair-api/src/entity"
	logs "github.com/marcovargas74/m74-fair-api/src/infrastructure/logs"
)

func TestServerAPI(t *testing.T) {

	logs.Start(false, "./fairAPI.log")

	assert.Equal(t, 1, 1)

}

func TestServerAPIDefault(t *testing.T) {

	tests := []struct {
		give      string
		wantValue string
		inData    string
	}{
		{
			give:      "Nobody return result",
			wantValue: "Endpoint not found",
			inData:    "Nobody",
		},
	}

	server := NewServerAPIMemory()
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodGet, "/", nil)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)

			assert.Equal(t, answer.Code, http.StatusAccepted)

		})
	}

}

func TestServerAPIDefaultPost(t *testing.T) {

	tests := []struct {
		give      string
		wantValue int
	}{
		{
			give:      "Default POST Endpoint test",
			wantValue: 405,
		},
	}

	server := NewServerAPIMemory()
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodPost, "/", nil)
			request.Header.Set("User-Agent", entity.UserAgentTest)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)

		})

	}

}

func TestCreateFairFail(t *testing.T) {

	tests := []struct {
		give      string
		wantValue int
		inData    string
	}{
		{
			give:      "Fair Endpoint test with NOBODY",
			wantValue: 405,
			inData:    "Nobody",
		},
		{
			give:      "Fair Endpoint test with empty char",
			wantValue: 404,
			inData:    "",
		},
		{
			give:      "Fair Endpoint test with Invalid data",
			wantValue: 405,
			inData:    "nome:feira",
		},
	}

	server := NewServerAPIMemory()
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request := entity.NewReqEndpointsPOST("/fair", tt.inData)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)
		})

	}
}

func TestCreateFair(t *testing.T) {

	tests := []struct {
		give      string
		wantValue int
		inData    string
	}{
		{
			give:      "Fair Endpoint test with NOBODY",
			wantValue: 405,
			inData:    "Nobody",
		},
	}

	server := NewServerAPIMemory()
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request := entity.NewReqEndpointsBodyPOST("/fair", tt.inData, tt.inData)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)
		})

	}
}

/*
func TestServerAPIStatusGet(t *testing.T) {

	tests := []struct {
		give      string
		wantValue string
	}{
		{
			give:      "status Endpoint test GET",
			wantValue: "{\"num_total_query\":0,\"start_time\":\"0001-01-01T00:00:00Z\",\"up_time\":9223372036.854776}",
		},
	}

	server := NewServerValidator("dev")
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodGet, "/status", nil)
			request.Header.Set("User-Agent", UserAgentTest)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, http.StatusOK)

		})

	}

}

func TestServerAPIStatusPost(t *testing.T) {

	tests := []struct {
		give      string
		wantValue intm74validatorapi "status Endpoint test POST",
			wantValue: 405,
		},
	}

	server := NewServerValidator("dev")
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodPost, "/status", nil)
			request.Header.Set("User-Agent", UserAgentTest)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)

		})

	}

}

func TestServerAPIQueryAllPost(t *testing.T) {

	tests := []struct {
		give      string
		wantValue int
	}{
		{
			give:      "All Query Endpoint test POST",
			wantValue: 405,
		},
	}

	server := NewServerValidator("dev")
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request, _ := http.NewRequest(http.MethodPost, "/all", nil)
			request.Header.Set("User-Agent", UserAgentTest)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)

		})

	}

}

func TestCallbackCpfsGET(t *testing.T) {

	tests := []struct {
		give      string
		wantValue int
		inData    string
	}{
		{
			give:      "cpfs Endpoint test with empty char",
			wantValue: 404,
			inData:    "",
		},
		{
			give:      "cpfs Endpoint test with NOBODY",
			wantValue: 400,
			inData:    "Nobody",
		},
		{
			give:      "cpfs Endpoint test with cnpj",
			wantValue: 404,
			inData:    "36.562.098/0001-18",
		},
	}

	server := NewServerValidator("dev")
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request := newReqEndpointsGET("/cpfs", tt.inData)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)
		})

	}

}

*/
