package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

/*
	{
	  "address": "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
	  "executionResult": {
	    "gasUsed": 21678,
	    "excepted": "None",
	    "newAddress": "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
	    "output": "0000000000000000000000000000000000000000000000000000000000000001",
	    "codeDeposit": 0,
	    "gasRefunded": 0,
	    "depositSize": 0,
	    "gasForDeposit": 0
	  },
	  "transactionReceipt": {
	    "stateRoot": "d44fc5ad43bae52f01ff7eb4a7bba904ee52aea6c41f337aa29754e57c73fba6",
	    "gasUsed": 21678,
	    "bloom": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	    "log": []
	  }
	}
*/
func TestEthCallRequest(t *testing.T) {
	//prepare request
	requestID, err := json.Marshal(1)
	request := eth.CallRequest{
		From: "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		To:   "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
	}
	requestRaw, err := json.Marshal(&request)
	if err != nil {
		panic(err)
	}
	requestParamsArray := []json.RawMessage{requestRaw}
	requestParamsArrayRaw, err := json.Marshal(requestParamsArray)
	if err != nil {
		panic(err)
	}

	requestRPC := &eth.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_protocolVersion",
		ID:      requestID,
		Params:  requestParamsArrayRaw,
	}

	d := doerMappedMock{make(map[string][]byte)}
	qtumClient, err := createMockedClient(d)

	//preparing response
	callContractResponse := qtum.CallContractResponse{
		Address: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		ExecutionResult: struct {
			GasUsed       int    "json:\"gasUsed\""
			Excepted      string "json:\"excepted\""
			NewAddress    string "json:\"newAddress\""
			Output        string "json:\"output\""
			CodeDeposit   int    "json:\"codeDeposit\""
			GasRefunded   int    "json:\"gasRefunded\""
			DepositSize   int    "json:\"depositSize\""
			GasForDeposit int    "json:\"gasForDeposit\""
		}{
			GasUsed:    21678,
			Excepted:   "None",
			NewAddress: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
			Output:     "0000000000000000000000000000000000000000000000000000000000000001",
		},
		TransactionReceipt: struct {
			StateRoot string        `json:"stateRoot"`
			GasUsed   int           `json:"gasUsed"`
			Bloom     string        `json:"bloom"`
			Log       []interface{} `json:"log"`
		}{
			StateRoot: "d44fc5ad43bae52f01ff7eb4a7bba904ee52aea6c41f337aa29754e57c73fba6",
			GasUsed:   21678,
			Bloom:     "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}
	err = d.AddResponse(1, qtum.MethodCallContract, callContractResponse)
	if err != nil {
		panic(err)
	}

	fromHexAddressResponse := qtum.FromHexAddressResponse("0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960")
	err = d.AddResponse(2, qtum.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		panic(err)
	}

	//preparing proxy & executing
	proxyEth := ProxyETHCall{qtumClient}
	if err != nil {
		panic(err)
	}

	got, err := proxyEth.Request(requestRPC)
	if err != nil {
		panic(err)
	}

	want := eth.CallResponse("0x0000000000000000000000000000000000000000000000000000000000000001")
	if !reflect.DeepEqual(got, &want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			requestRPC,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}
