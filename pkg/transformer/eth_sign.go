package transformer

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHGetLogs implements ETHProxy
type ProxyETHSign struct {
	*qtum.Qtum
}

func (p *ProxyETHSign) Method() string {
	return "eth_sign"
}

func (p *ProxyETHSign) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.SignRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	base58Addr, err := p.FromHexAddress(utils.RemoveHexPrefix(req.Account))
	if err != nil {
		return nil, err
	}

	sig, err := p.SignMessage(base58Addr, string(req.Message))
	if err != nil {
		return nil, err
	}

	// base64 -> hex

	sigdata, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return nil, err
	}

	return eth.SignResponse("0x" + hex.EncodeToString(sigdata)), nil
}
