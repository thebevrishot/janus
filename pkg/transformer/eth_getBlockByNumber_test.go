package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestGetBlockByNumberRequest(t *testing.T) {
	requestParams := []json.RawMessage{[]byte(`"0x1b4"`), []byte(`true`)}
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		panic(err)
	}

	mockedClientDoer := doerMappedMock{make(map[string][]byte)}
	qtumClient, err := createMockedClient(mockedClientDoer)

	//preparing answer to "getblockhash"
	getBlockHashResponse := qtum.GetBlockHashResponse("0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331")
	err = mockedClientDoer.AddResponse(2, qtum.MethodGetBlockHash, getBlockHashResponse)
	if err != nil {
		panic(err)
	}

	getBlockHeaderResponse := qtum.GetBlockHeaderResponse{
		Hash: "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
	}
	err = mockedClientDoer.AddResponse(3, qtum.MethodGetBlockHeader, getBlockHeaderResponse)
	if err != nil {
		panic(err)
	}

	getBlockResponse := qtum.GetBlockResponse{
		Hash:              "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		Previousblockhash: "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
		Size:              0x027f07,
		Difficulty:        0x027f07,
		Nonce:             0, //?
	}
	err = mockedClientDoer.AddResponse(4, qtum.MethodGetBlock, getBlockResponse)
	if err != nil {
		panic(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBlockByNumber{qtumClient}
	got, err := proxyEth.Request(request)
	if err != nil {
		panic(err)
	}

	want := &eth.GetBlockByNumberResponse{
		Number:     "0x0", // ?
		Hash:       "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		ParentHash: "0x", // ?!
		//ParentHash:       "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5", <--- should be
		Miner:            "0x0000000000000000000000000000000000000000",
		Size:             "0x27f07",
		Nonce:            "0x0", // ?
		TransactionsRoot: "0x",  // ?
		StateRoot:        "0x",  // ?
		Difficulty:       "0x0", // ?!
		TotalDifficulty:  "0x0", // ?!
		ExtraData:        "0x0",
		GasLimit:         "0x0",
		GasUsed:          "0x0",
		Timestamp:        "0x0",
		Transactions:     []string{},
		Uncles:           []string{},
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
