package transformer

import (
	"log"

	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	//"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHSendRawTransaction implements ETHProxy
type ProxyETHSendRawTransaction struct {
	*qtum.Qtum
}

func (p *ProxyETHSendRawTransaction) Method() string {
	return "eth_sendRawTransaction"
}

func (p *ProxyETHSendRawTransaction) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req string
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	if p.Chain() == qtum.ChainRegTest {
		defer func() {
			if _, generateErr := p.Generate(1, nil); generateErr != nil {
				log.Println("generate block err: ", generateErr)
			}
		}()
	}

	return nil, errors.New("Unknown operation")
}
