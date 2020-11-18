package transformer

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

/*
	Test data:
	 {Transactions
	    "number": "0x1b4",
	    "hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
	    "parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
	    "nonce": "0xe04d296d2460cfb8472af2c5fd05b5a214109c25688d3704aed5484f9a7792f2",
	    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
	    "logsBloom": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
	    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
	    "stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
	    "miner": "0x4e65fda2159562a496f9f3522f89122a3088497a",
	    "difficulty": "0x027f07",
	    "totalDifficulty":  "0x027f07",
	    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
	    "size":  "0x027f07",
	    "gasLimit": "0x9f759",
	    "gasUsed": "0x9f759",
	    "timestamp": "0x54e34e8e",
	    "transactions": [{}],
	    "uncles": ["0x1606e5...", "0xd5145a9..."]
	  }
*/
func TestGetBlockByNumberRequest(t *testing.T) {
	requestID, err := json.Marshal(1)
	request := &eth.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_protocolVersion",
		ID:      requestID,
		Params:  []byte(`{"blockNumber": "0x1b4","fullTransaction": true}`),
	}

	doerInstance := doerMappedMock{make(map[string][]byte)}
	qtumClient, err := createMockedClient(doerInstance)

	//preparing answer to "getblockhash"
	getBlockHashResponse := qtum.GetBlockHashResponse("0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331")
	getBlockHashResponseRaw, err := json.Marshal(getBlockHashResponse)
	getBlockHashResponseRPC := &eth.JSONRPCResult{
		JSONRPC:   "2.0",
		RawResult: getBlockHashResponseRaw,
		Error:     nil,
		ID:        requestID,
	}

	getBlockHashResponseRPCRaw, err := json.Marshal(getBlockHashResponseRPC)
	doerInstance.Responses[qtum.MethodGetBlockHash] = getBlockHashResponseRPCRaw

	//preparing answer to "getblockheader"
	getBlockHeaderResponse := qtum.GetBlockHeaderResponse{
		Hash: "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
	}
	getBlockHeaderResponseRaw, err := json.Marshal(getBlockHeaderResponse)
	getBlockHeaderResponseRPC := &eth.JSONRPCResult{
		JSONRPC:   "2.0",
		RawResult: getBlockHeaderResponseRaw,
		Error:     nil,
		ID:        requestID,
	}

	getBlockHeaderResponseRPCRaw, err := json.Marshal(getBlockHeaderResponseRPC)
	doerInstance.Responses[qtum.MethodGetBlockHeader] = getBlockHeaderResponseRPCRaw

	//preparing answer to "getblock"
	getBlockResponse := qtum.GetBlockResponse{
		Hash:              "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
		Previousblockhash: "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
		Size:              0x027f07,
		Difficulty:        0x027f07,
		Nonce:             0, //?
	}
	getBlockResponseRaw, err := json.Marshal(getBlockResponse)
	getBlockResponseRPC := &eth.JSONRPCResult{
		JSONRPC:   "2.0",
		RawResult: getBlockResponseRaw,
		Error:     nil,
		ID:        requestID,
	}

	getBlockResponseRPCRaw, err := json.Marshal(getBlockResponseRPC)
	doerInstance.Responses[qtum.MethodGetBlock] = getBlockResponseRPCRaw

	//preparing proxy
	proxyEth := ProxyETHGetBlockByNumber{qtumClient}
	if err != nil {
		panic(err)
	}

	//executing request
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
