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

			request := entity.NewReqEndpointsPOST("/fairs", tt.inData)
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

			request := entity.NewReqEndpointsBodyPOST("/fairs", tt.inData, tt.inData)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)
		})

	}
}

func TestCreateUpdateFail(t *testing.T) {

	tests := []struct {
		give      string
		wantValue int
		inData    string
	}{
		{
			give:      "Fair Endpoint test with empty char",
			wantValue: 404,
			inData:    "",
		},
		{
			give:      "Fair Endpoint test with iNVALID id",
			wantValue: 400,
			inData:    "123",
		},
		{
			give:      "Fair Endpoint test with iNVALID id",
			wantValue: 404,
			inData:    "4f45e8f1-c75d-4330-ba7b-58c691d61104",
		},
	}

	server := NewServerAPIMemory()
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {

			request := entity.NewReqEndpointsPUT("/fairs", tt.inData)
			answer := httptest.NewRecorder()

			server.ServeHTTP(answer, request)
			assert.Equal(t, answer.Code, tt.wantValue)
		})

	}
}
