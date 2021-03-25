package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/qtumproject/janus/pkg/eth"
)

type ProxyETHChainId struct{}

func (p *ProxyETHChainId) Method() string {
	return "eth_chainId"
}

func (p *ProxyETHChainId) Request(req *eth.JSONRPCRequest) (interface{}, error) {
	var chainId = big.NewInt(81)

	return eth.ChainIdResponse(hexutil.EncodeBig(chainId)), nil
}
