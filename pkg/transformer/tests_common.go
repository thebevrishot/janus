package transformer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

//copy of qtum.doer interface
type doer interface {
	Do(*http.Request) (*http.Response, error)
}

//type for mocking requests to client for simple requests
type doerMock struct {
	Response []byte
}

func (d doerMock) Do(*http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(d.Response)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

//type for mocking requests to client with request -> responce mapping
type doerMappedMock struct {
	Responses map[string][]byte
}

func (d doerMappedMock) Do(request *http.Request) (*http.Response, error) {
	requestJSON, err := parseRequestFromBody(request)
	if err != nil {
		return nil, err
	}

	responseWriter := ioutil.NopCloser(bytes.NewReader(d.Responses[requestJSON.Method]))
	return &http.Response{
		StatusCode: 200,
		Body:       responseWriter,
	}, nil
}

func prepareRawResponse(requestID int, responseResult interface{}) ([]byte, error) {
	requestIDRaw, err := json.Marshal(requestID)
	if err != nil {
		return nil, err
	}

	responseResultRaw, err := json.Marshal(responseResult)
	if err != nil {
		return nil, err
	}

	responseRPC := &eth.JSONRPCResult{
		JSONRPC:   "2.0",
		RawResult: responseResultRaw,
		Error:     nil,
		ID:        requestIDRaw,
	}

	responseRPCRaw, err := json.Marshal(responseRPC)

	return responseRPCRaw, err
}

func (d *doerMappedMock) AddResponse(requestID int, requestType string, responseResult interface{}) error {
	responseRaw, err := prepareRawResponse(requestID, responseResult)
	if err != nil {
		return err
	}

	d.Responses[requestType] = responseRaw
	return nil
}

func parseRequestFromBody(request *http.Request) (*eth.JSONRPCRequest, error) {
	requestJSON := eth.JSONRPCRequest{}
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(requestBody, &requestJSON)
	if err != nil {
		return nil, err
	}

	return &requestJSON, err
}

func createMockedClient(doerInstance doer) (qtumClient *qtum.Qtum, err error) {
	qtumJSONRPC, err := qtum.NewClient(true, "http://user:pass@mocked", qtum.SetDoer(doerInstance))
	if err != nil {
		return
	}

	qtumClient, err = qtum.New(qtumJSONRPC, "test")
	return
}

func mustMarshalIndent(v interface{}, prefix, indent string) []byte {
	res, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return res
}
