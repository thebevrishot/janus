package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestGetTransactionReceiptForNonVMTransaction(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{[]byte(`"0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5"`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing client response
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodGetTransactionReceipt, []byte("[]"))
	if err != nil {
		t.Fatal(err)
	}

	rawTransactionResponse := &qtum.GetRawTransactionResponse{
		BlockHash: internal.GetTransactionByHashBlockHash,
	}
	err = mockedClientDoer.AddResponseWithRequestID(3, qtum.MethodGetRawTransaction, rawTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}

	err = mockedClientDoer.AddResponseWithRequestID(4, qtum.MethodGetBlock, internal.GetBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetTransactionReceipt{qtumClient}
	got, err := proxyEth.Request(request, nil)
	if err != nil {
		t.Fatal(err)
	}

	want := eth.GetTransactionReceiptResponse{
		TransactionHash:   "0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		TransactionIndex:  "0x1",
		BlockHash:         "0xbba11e1bacc69ba535d478cf1f2e542da3735a517b0b8eebaf7e6bb25eeb48c5",
		BlockNumber:       "0xf8f",
		GasUsed:           "0x0",
		CumulativeGasUsed: "0x0",
		LogsBloom:         eth.EmptyLogsBloom,
		Status:            "0x0",
	}
	if !reflect.DeepEqual(got, &want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			want,
			got,
		)
	}
}
