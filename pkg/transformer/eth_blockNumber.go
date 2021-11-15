package transformer

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHBlockNumber implements ETHProxy
type ProxyETHBlockNumber struct {
	*qtum.Qtum
	cacher *BlockSyncer
}

func (p *ProxyETHBlockNumber) Method() string {
	return "eth_blockNumber"
}

func (p *ProxyETHBlockNumber) WithBlockCacher(cacher *BlockSyncer) *ProxyETHBlockNumber {
	p.cacher = cacher
	return p
}

func (p *ProxyETHBlockNumber) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	return p.request(c, 5)
}

func (p *ProxyETHBlockNumber) request(c echo.Context, retries int) (*eth.BlockNumberResponse, error) {

	if p.cacher != nil {
		block, ok := p.cacher.GetLatestBlock()
		if ok && block != nil {
			return (*eth.BlockNumberResponse)(&block.Number), nil
		}
	}

	qtumresp, err := p.Qtum.GetBlockCount()
	if err != nil {
		if retries > 0 && strings.Contains(err.Error(), qtum.ErrTryAgain.Error()) {
			ctx := c.Request().Context()
			t := time.NewTimer(500 * time.Millisecond)
			select {
			case <-ctx.Done():
				return nil, err
			case <-t.C:
				// fallthrough
			}
			return p.request(c, retries-1)
		}
		return nil, err
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp), nil
}

func (p *ProxyETHBlockNumber) ToResponse(qtumresp *qtum.GetBlockCountResponse) *eth.BlockNumberResponse {
	hexVal := hexutil.EncodeBig(qtumresp.Int)
	ethresp := eth.BlockNumberResponse(hexVal)
	return &ethresp
}
