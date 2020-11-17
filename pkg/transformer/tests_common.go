package transformer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/qtumproject/janus/pkg/qtum"
)

//type for mocking requests to client
type doerMock struct {
	response []byte
}

//implementation of qtum.Doer interface
func (d doerMock) Do(*http.Request) (*http.Response, error) {
	r := ioutil.NopCloser(bytes.NewReader([]byte(d.response)))
	return &http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil
}

//creates instance of ProxyETHGetTransactionByHash with mocked http-response
func createMockedClient(response []byte) (qtumClient *qtum.Qtum, err error) {
	doer := doerMock{response}
	qtumJSONRPC, err := qtum.NewClient(true, "http://user:pass@mocked", qtum.SetDoer(doer), qtum.SetDebug(true))
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
