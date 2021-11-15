package transformer

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// ProxyETHGetBlockByNumber implements ETHProxy
type ProxyETHGetBlockByNumber struct {
	*qtum.Qtum
	cacher *BlockSyncer
}

func (p *ProxyETHGetBlockByNumber) Method() string {
	return "eth_getBlockByNumber"
}

func (p *ProxyETHGetBlockByNumber) Request(rpcReq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	req := new(eth.GetBlockByNumberRequest)
	if err := unmarshalRequest(rpcReq.Params, req); err != nil {
		return nil, errors.WithMessage(err, "couldn't unmarhsal rpc request")
	}
	return p.request(req)
}

func (p *ProxyETHGetBlockByNumber) WithBlockPoller(cacher *BlockSyncer) *ProxyETHGetBlockByNumber {
	p.cacher = cacher
	return p
}

func (p *ProxyETHGetBlockByNumber) request(req *eth.GetBlockByNumberRequest) (*eth.GetBlockByNumberResponse, error) {
	if p.cacher != nil && !req.FullTransaction {
		block, ok := p.cacher.GetBlock(req.BlockNumber)
		if ok {
			if block != nil {
				return block, nil
			} else {
				var blockReq string
				if err := json.Unmarshal(req.BlockNumber, &blockReq); err != nil {
					fmt.Println("fail to unmarshal", err)
				}
				return nil, errors.New("couldn't get block number by parameter " + blockReq)
			}
		}
	}

	blockNum, err := getBlockNumberByRawParam(p.Qtum, req.BlockNumber, false)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get block number by parameter")
	}

	blockHash, err := proxyETHGetBlockByHash(p, p.Qtum, blockNum)
	if err != nil {
		return nil, err
	}
	if blockHash == nil {
		return nil, nil
	}

	var (
		getBlockByHashReq = &eth.GetBlockByHashRequest{
			BlockHash:       string(*blockHash),
			FullTransaction: req.FullTransaction,
		}
		proxy = &ProxyETHGetBlockByHash{Qtum: p.Qtum}
	)
	block, err := proxy.request(getBlockByHashReq)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get block by hash")
	}
	if blockNum != nil {
		p.GetDebugLogger().Log("function", p.Method(), "request", string(req.BlockNumber), "msg", "Successfully got block by number", "result", blockNum.String())
	}
	return block, nil
}

// Properly handle unknown blocks
func proxyETHGetBlockByHash(p ETHProxy, q *qtum.Qtum, blockNum *big.Int) (*qtum.GetBlockHashResponse, error) {
	resp, err := q.GetBlockHash(blockNum)
	if err != nil {
		if err == qtum.ErrInvalidParameter {
			// block doesn't exist, ETH rpc returns null
			/**
			{
				"jsonrpc": "2.0",
				"id": 1234,
				"result": null
			}
			**/
			q.GetDebugLogger().Log("function", p.Method(), "request", blockNum.String(), "msg", "Unknown block")
			return nil, nil
		}
		return nil, errors.WithMessage(err, "couldn't get block hash")
	}
	return &resp, err
}
