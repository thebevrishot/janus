package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestGetTransactionByHashRequest(t *testing.T) {
	//prepare request
	id, err := json.Marshal(int64(67))
	request := &eth.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_protocolVersion",
		ID:      id,
		Params:  json.RawMessage(`["0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5"]`),
	}
	if err != nil {
		panic(err)
	}

	//prepare expected response & client
	mockedResponseResult := qtum.GetTransactionResponse{
		Amount:            0,
		Fee:               -0.2012,
		Confirmations:     2,
		Blockhash:         "ea26fd59a2145dcecd0e2f81b701019b51ca754b6c782114825798973d8187d6",
		Blockindex:        2,
		Blocktime:         1533092896,
		Txid:              "11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		Time:              1533092879,
		Timereceived:      1533092879,
		Bip125Replaceable: "no",
		Details: []*qtum.TransactionDetail{{Account: "",
			Category:  "send",
			Amount:    0,
			Vout:      0,
			Fee:       -0.2012,
			Abandoned: false}},
		Hex: "020000000159c0514feea50f915854d9ec45bc6458bb14419c78b17e7be3f7fd5f563475b5010000006a473044022072d64a1f4ea2d54b7b05050fc853ab192c91cc5ca17e23007867f92f2ab59d9202202b8c9ab9348c8edbb3b98b1788382c8f37642ec9bd6a4429817ab79927319200012103520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140feffffff02000000000000000063010403400d0301644440c10f190000000000000000000000006b22910b1e302cf74803ffd1691c2ecb858d3712000000000000000000000000000000000000000000000000000000000000000a14be528c8378ff082e4ba43cb1baa363dbf3f577bfc260e66272970100001976a9146b22910b1e302cf74803ffd1691c2ecb858d371288acb00f0000",
	}
	mockedResponseResultRaw, err := json.Marshal(mockedResponseResult)
	mockedResponse := &eth.JSONRPCResult{
		JSONRPC:   "2.0",
		RawResult: mockedResponseResultRaw,
		Error:     nil,
		ID:        id,
	}

	responseRaw, err := json.Marshal(mockedResponse)
	if err != nil {
		panic(err)
	}
	doer := doerMock{responseRaw}
	qtumClient, err := createMockedClient(doer)
	proxyEth := ProxyETHGetTransactionByHash{qtumClient}
	if err != nil {
		panic(err)
	}

	//execute request
	got, err := proxyEth.Request(request)
	if err != nil {
		panic(err)
	}

	want := &eth.GetTransactionByHashResponse{
		Hash:      "0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		BlockHash: "0xea26fd59a2145dcecd0e2f81b701019b51ca754b6c782114825798973d8187d6",
		Value:     "0x0",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"error\ninput: %s\nwant: %s\ngot: %s",
			request,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}
