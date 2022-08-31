package entity

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	erroMsg       = "Test fail got Value[%v], wait Value [%v]"
	UserAgentTest = "self_test"
)

//NewReqEndpointsPOST Genetic POST endpoint to test
func NewReqEndpointsPOST(urlPrefix, urlName string) *http.Request {
	request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", urlPrefix, urlName), nil)
	if error != nil {
		panic(error)
	}

	request.Header.Set("User-Agent", UserAgentTest)
	return request
}

//NewReqEndpointsPUT Genetic PUT endpoint to test
func NewReqEndpointsPUT(urlPrefix, urlName string) *http.Request {
	request, error := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", urlPrefix, urlName), nil)
	if error != nil {
		panic(error)
	}

	request.Header.Set("User-Agent", UserAgentTest)
	return request
}

//NewReqEndpointsBodyPOST Genetic POST endpoint to test
func NewReqEndpointsBodyPOST(urlPrefix, urlName string, json string) *http.Request {

	jsonBody := []byte(`{
		"name": "VILA TESTE",
		"District": "dstrito teste",
		"Region5": "test",
		"neighborhood" :"bairo"
	}
      `)

	bodyReader := bytes.NewReader(jsonBody)

	request, error := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", urlPrefix, urlName), bodyReader)
	if error != nil {
		panic(error)
	}

	request.Header.Set("User-Agent", UserAgentTest)
	return request
}
