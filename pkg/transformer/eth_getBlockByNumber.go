package transformer

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/labstack/echo"
	"github.com/patrickmn/go-cache"
	_ "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

const (
	DefaultBlockCache = 128
)

// ProxyETHGetBlockByNumber implements ETHProxy
type ProxyETHGetBlockByNumber struct {
	*qtum.Qtum
	blockCache         *lru.Cache
	blockHashCacheLock sync.RWMutex
	blockHashCache     *cache.Cache
	enableCache        bool
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

func (p *ProxyETHGetBlockByNumber) WithCache() *ProxyETHGetBlockByNumber {
	p.blockCache, _ = lru.New(DefaultBlockCache)
	p.blockHashCacheLock = sync.RWMutex{}
	p.blockHashCache = cache.New(2*time.Second, 10*time.Second)
	p.enableCache = true
	return p
}

func (p *ProxyETHGetBlockByNumber) getBlockHash(blockNumber json.RawMessage) (*qtum.GetBlockHashResponse, error) {
	// Look up in cache
	waitForHash := false
	if p.enableCache {
		p.blockHashCacheLock.Lock()
		defer p.blockHashCacheLock.Unlock()
		val, found := p.blockHashCache.Get(string(blockNumber))

		if !found {
			// set to let other know someone query hash for them
			if err := p.blockHashCache.Add(string(blockNumber), nil, 5*time.Second); err != nil {
				return nil, err
			}
		} else if hashCandidate, isHash := val.(*qtum.GetBlockHashResponse); isHash {
			return hashCandidate, nil
		} else {
			waitForHash = true
		}
	}

	// wait for other to query hash
	if waitForHash {
		timeout := time.NewTicker(60 * time.Second)
		interval := time.NewTicker(200 * time.Millisecond)

	Outer:
		for {
			select {
			case <-timeout.C:
				break Outer
			case <-interval.C:
				// check if hash is set
				val, _ := p.blockHashCache.Get(string(blockNumber))
				if hashCandidate, isHash := val.(*qtum.GetBlockHashResponse); isHash {
					fmt.Println("logfoo: Hit query hash")
					return hashCandidate, nil
				}
			}
		}
	}

	blockNum, err := getBlockNumberByRawParam(p.Qtum, blockNumber, false)
	if err != nil {
		return nil, errors.WithMessage(err, "couldn't get block number by parameter")
	}

	blockHash, err := proxyETHGetBlockByHash(p, p.Qtum, blockNum)
	if err != nil {
		return nil, err
	}

	if p.enableCache {
		p.blockHashCache.Add(string(blockNumber), blockHash, 2*time.Second)
	}

	return blockHash, nil
}

func (p *ProxyETHGetBlockByNumber) request(req *eth.GetBlockByNumberRequest) (*eth.GetBlockByNumberResponse, error) {
	start := time.Now()
	blockHash, err := p.getBlockHash(req.BlockNumber)
	if blockHash == nil {
		return nil, err
	}
	fmt.Printf("logfoo: Query blockhash: %v\n", time.Since(start))

	if ok, _ := p.blockCache.ContainsOrAdd(*blockHash, nil); ok {

		// If contain block then return
		val, _ := p.blockCache.Get(*blockHash)
		if block, ok := val.(*eth.GetBlockByNumberResponse); ok {
			fmt.Println("logfoo: Hit", string(*blockHash))
			return block, nil
		}

		// Already has or someone querying it
		timeOut := time.NewTicker(time.Second * 10)
		interval := time.NewTicker(time.Millisecond * 100)

	OuterLoop:
		for {
			select {
			case <-timeOut.C:
				break OuterLoop
			case <-interval.C:
				val, _ := p.blockCache.Get(*blockHash)
				if _, ok := val.(*eth.GetBlockByNumberResponse); ok {
					break OuterLoop
				}
			}
		}

		// If found block then return
		val, _ = p.blockCache.Get(*blockHash)
		if block, ok := val.(*eth.GetBlockByNumberResponse); ok {
			fmt.Println("logfoo: Hit", string(*blockHash))
			return block, nil
		}

		// Block was not queried in time
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
	// if blockNum != nil {
	// 	p.GetDebugLogger().Log("function", p.Method(), "request", string(req.BlockNumber), "msg", "Successfully got block by number", "result", blockNum.String())
	// }

	p.blockCache.Add(*blockHash, block)

	fmt.Println("logfoo: Stored", string(*blockHash))
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
