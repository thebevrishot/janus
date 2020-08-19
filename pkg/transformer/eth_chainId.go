package transformer

import (
	"github.com/qtumproject/janus/pkg/eth"
)

// curl -X POST --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'
// // Result
// {
//   "id": 83,
//   "jsonrpc": "2.0",
//   "result": "0x3d" // 61
// }

// ProxyETHChainID implements ProxyETHChainID
type ProxyETHChainID struct {
}

func (p *ProxyETHChainID) Method() string {
	return "eth_chainId"
}

func (p *ProxyETHChainID) Request(_ *eth.JSONRPCRequest) (interface{}, error) {
	return p.request()
}

func (p *ProxyETHChainID) request() (string, error) {

	networkID := "0x1024"
	// switch qtumresp.Chain {
	// case "regtest":
	// 	// See: https://github.com/trufflesuite/ganache/issues/112 for an idea on how to generate an ID.
	// 	// https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version
	// 	networkID = "0x1024"
	// default:
	// 	networkID = qtumresp.Chain
	// }

	// resp := eth.ETHChainIDResponse(networkID)
	return networkID, nil
}
