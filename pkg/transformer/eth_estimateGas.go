package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHEstimateGas struct {
	*ProxyETHCall
}

func (p *ProxyETHEstimateGas) Method() string {
	return "eth_estimateGas"
}

func (p *ProxyETHEstimateGas) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var ethreq eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &ethreq); err != nil {
		return nil, err
	}

	// When deploying a contract, the To address is empty. And we won't be
	// able to get an gas estimate from callcontract.
	if ethreq.To == "" {
		// Just return 10 qtum
		return eth.EstimateGasResponse(hexutil.EncodeUint64(uint64(10 * 1e9))), nil
	}

	// eth req -> qtum req
	qtumreq, err := p.ToRequest(&ethreq)
	if err != nil {
		return nil, err
	}

	// qtum [code: -5] Incorrect address occurs here
	qtumresp, err := p.CallContract(qtumreq)
	if err != nil {
		return nil, err
	}

	return p.toResp(qtumresp)
}

func (p *ProxyETHEstimateGas) toResp(qtumresp *qtum.CallContractResponse) (*eth.EstimateGasResponse, error) {
	gas := eth.EstimateGasResponse(hexutil.EncodeUint64(uint64(qtumresp.ExecutionResult.GasUsed)))
	return &gas, nil
}
