package transformer

import (
	"fmt"
	"math/big"
	"regexp"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// ProxyETHCall implements ETHProxy
type ProxyETHCall struct {
	*qtum.Qtum
}

func (p *ProxyETHCall) Method() string {
	return "eth_call"
}

func (p *ProxyETHCall) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	return p.request(&req)
}

func (p *ProxyETHCall) request(ethreq *eth.CallRequest) (interface{}, error) {
	// eth req -> qtum req
	qtumreq, err := p.ToRequest(ethreq)
	if err != nil {
		return nil, err
	}

	qtumresp, err := p.CallContract(qtumreq)
	if err != nil {
		return nil, err
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp), nil
}

func (p *ProxyETHCall) ToRequest(ethreq *eth.CallRequest) (*qtum.CallContractRequest, error) {
	from := ethreq.From
	var err error
	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, err
		}
	}

	return &qtum.CallContractRequest{
		To:   ethreq.To,
		From: from,
		Data: ethreq.Data,
		// TODO: qtum [code: -3] Invalid value for gasLimit (Minimum is: 10000)
		// Incorrect gas format
		GasLimit: big.NewInt(10000),
		//GasLimit: ethreq.Gas.Int,
	}, nil
}

func (p *ProxyETHCall) ToResponse(qresp *qtum.CallContractResponse) interface{} {
	excepted := qresp.ExecutionResult.Excepted
	data := utils.AddHexPrefix(qresp.ExecutionResult.Output)
	if excepted != "None" {
		//Get transaction: Might not be needed
		//get code -> data
		//decode message
		
		message := DecodeMessage(data)

		


		/*
		 	code = eth_call({transaction}, blockNumber)
			message = decodeMessage(code)
			return eth.JSONRPCError{
				"message": "VM Exception while processing transaction: revert CapperRole: caller does not have the Capper role",
				"code": -32000,
				"data": {
				"0x90f48d0a98bf9e18bba07a9538da7a1c21fcccab4aaca3afd60b6d0677098117": {
					"error": "revert",
					"program_counter": 621,
					"return": "0x08c379a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000030436170706572526f6c653a2063616c6c657220646f6573206e6f742068617665207468652043617070657220726f6c6500000000000000000000000000000000",
					"reason": "CapperRole: caller does not have the Capper role"
				},
				"stack": "o: VM Exception while processing transaction: revert CapperRole: caller does not have the Capper role\n    at Function.o.fromResults (/Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:10:81931)\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:47:121235\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:61:1853218\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:26124\n    at i (/Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:41179)\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:61:1210468\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:105466\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:32:392\n    at c (/Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:32:5407)\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:32:317\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:61:1871644\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:23237\n    at o (/Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:26646)\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:26124\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:61:1864681\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:61:1862544\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:61:1889757\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:23237\n    at o (/Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:26646)\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:26124\n    at /Users/howard/src/openzeppelin-solidity/node_modules/ganache-cli/build/ganache-core.node.cli.js:2:5439\n    at FSReqCallback.args [as oncomplete] (fs.js:145:20)",
				"name": "o"
				}
			}
		*/
		return &eth.JSONRPCError{
			Code:    -32000,
			Message: fmt.Sprintf("VM exception: %s", excepted),
			// To see how eth_call supports revert reason, see:
			// https://gist.github.com/hayeah/795bc18a683053218fb3ff5032d31144
			//
			// Data: ...
		}
	}

	
	qtumresp := eth.CallResponse(data)
	return &qtumresp
}

func (p *ProxyETHCall) DecodeMessage(code string) string {

	var codeString string

	const FnSelectorByteLength = 4
	const DataOffsetByteLength = 32
	const StrLengthByteLength = 32
	const StrLengthStartPos = 2 + ((FnSelectorByteLength + DataOffsetByteLength) * 2)
	const StrDataStartPos = 2 + ((FnSelectorByteLength + DataOffsetByteLength + StrLengthByteLength) * 2)

	re := regexp.MustCompile(`/0+$/`)
	codeString = "0x" + strings.ReplaceAllString(code[138:], "")

	// If the codeString is an odd number of characters, add a trailing 0
	if len(codeString) % 2 == 1 {
		codeString += "0"
	}
	

	return hexutil.Decode(codeString) //Have to handle decoding error
}
