package presenter

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestSelectKeySearch(t *testing.T) {

	tests := []struct {
		give      string
		inData    string
		wantKey   string
		wantValue string
	}{
		{
			give:      "Fair SelectKeySearch test with empty char",
			inData:    "",
			wantKey:   "",
			wantValue: "",
		},
		{
			give:      "Fair SelectKeySearch test with name key",
			inData:    "/fairs?name=vila",
			wantKey:   "name",
			wantValue: "vila",
		},
		{
			give:      "Fair SelectKeySearch test with district key",
			inData:    "/fairs?district=VILA FORMOSA",
			wantKey:   "district",
			wantValue: "VILA FORMOSA",
		},
		{
			give:      "Fair SelectKeySearch test with region5 key",
			inData:    "/fairs?region5=Sul",
			wantKey:   "region5",
			wantValue: "Sul",
		},
		{
			give:      "Fair SelectKeySearch test with neighborhood key",
			inData:    "/fairs?neighborhood=VILA FORMOSA",
			wantKey:   "neighborhood",
			wantValue: "VILA FORMOSA",
		},
		{
			give:      "Fair SelectKeySearch test with region5= key",
			inData:    "/fairs?region5= ",
			wantKey:   "region5",
			wantValue: " ",
		},
		{
			give:      "Fair SelectKeySearch test with name= key",
			inData:    "/fairs?name= ",
			wantKey:   "name",
			wantValue: " ",
		},
		{
			give:      "Fair SelectKeySearch test with name= fail key",
			inData:    "/fairs?name=",
			wantKey:   "",
			wantValue: "",
		}}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, tt.inData, nil)
			key, value := SelectKeySearch(request)

			fmt.Printf(".....testS[%s] [%s]\n", key, value)

			assert.Equal(t, key, tt.wantKey)
			assert.Equal(t, value, tt.wantValue)
		})

	}

}

/*
func TestServerAPIDefault(t *testing.T) {

	tests := []struct {toJSON := presenter.NewCreateFairPresenter(d)
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

/*
toJSON := presenter.NewCreateFairPresenter(d)
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
