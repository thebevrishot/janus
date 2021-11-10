package transformer

import (
	"container/list"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

type BlockPoller struct {
	Qtum   *qtum.Qtum
	lock   sync.RWMutex
	blocks *list.List
	synced bool
	limit  int
}

func (p *BlockPoller) pullBlock(blockNumber *big.Int) (*eth.GetBlockByHashResponse, error) {
	proxy := ProxyETHGetBlockByNumber{Qtum: p.Qtum}

	blockHash, err := proxyETHGetBlockByHash(&proxy, p.Qtum, blockNumber)
	if err != nil {
		return nil, err
	}

	var (
		getBlockByHashReq = &eth.GetBlockByHashRequest{
			BlockHash:       string(*blockHash),
			FullTransaction: false,
		}
		proxyETHGetBlockByHash = &ProxyETHGetBlockByHash{Qtum: p.Qtum}
	)
	return proxyETHGetBlockByHash.request(getBlockByHashReq)
}

func (p *BlockPoller) clearBlocks() {
	p.synced = false

	// Assume lock is locked
	itr := p.blocks.Front()

	// No valid block then cleanup
	for itr != nil {
		delItr := itr
		itr = itr.Next()
		p.blocks.Remove(delItr)
	}
}

func (p *BlockPoller) loopSync() error {

	for {
		// Query block count
		blockCountResp, err := p.Qtum.GetBlockCount()
		if err != nil {
			p.Qtum.GetLogger().Log("function", "loopSync", "message", "fail to query blockcount", "error", err)
			continue
		}
		upstreamBlock := blockCountResp.Int

		// Get local block
		blockItr := p.blocks.Back()

		var localBlock *big.Int
		var localHash = ""
		if blockItr != nil {
			if block, ok := blockItr.Value.(*eth.GetBlockByHashResponse); ok {
				localHash = block.Hash
				blockNum, _ := strconv.ParseInt(block.Number[2:], 16, 64)
				localBlock = big.NewInt(int64(blockNum))
			}
		} else {
			localBlock = big.NewInt(0).Sub(upstreamBlock, big.NewInt(1))
		}

		if localBlock.Cmp(upstreamBlock) < 0 {
			newBlock, err := p.pullBlock(big.NewInt(0).Add(localBlock, big.NewInt(1)))
			if err != nil {
				p.Qtum.GetLogger().Log("function", "loopSync", "message", "fail to query block", "error", err)
				continue
			}

			if localHash == "" || newBlock.ParentHash == localHash {
				p.lock.Lock()
				p.blocks.PushBack(newBlock)
				p.synced = true
				p.lock.Unlock()
			} else {
				p.Qtum.GetLogger().Log("function", "loopSync", "message", "last block is invalid", "error", err)
				p.lock.Lock()
				p.clearBlocks()
				p.lock.Unlock()
				continue
			}
		} else if localBlock.Cmp(upstreamBlock) > 0 {
			p.lock.Lock()
			p.clearBlocks()
			p.lock.Unlock()
			continue
		} else {
			upstreamHash, err := p.Qtum.GetBlockHash(localBlock)
			if err != nil {
				p.Qtum.GetLogger().Log("function", "loopSync", "message", "Fail to get block hash of local height", "error", err)
				continue
			}

			if "0x"+upstreamHash != qtum.GetBlockHashResponse(localHash) {
				p.Qtum.GetLogger().Log("function", "loopSync", "message", "Upstream hash is not match with local", "error", err)
				p.lock.Lock()
				p.clearBlocks()
				p.lock.Unlock()
				continue
			}

			// If last block is corrent just sleep
			time.Sleep(5 * time.Second)
		}

		// Cleanup old
		for p.blocks.Len() > p.limit {
			p.lock.Lock()
			p.blocks.Remove(p.blocks.Front())
			p.lock.Unlock()
		}
	}
}

func (p *BlockPoller) GetBlock(blockNumber json.RawMessage) (*eth.GetBlockByHashResponse, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.synced {
		return nil, false
	}

	var blockNumberStr string
	if err := json.Unmarshal(blockNumber, &blockNumberStr); err != nil {
		return nil, false
	}

	var element *list.Element
	switch blockNumberStr {
	case "latest":
		element = p.blocks.Back()
		if block, ok := element.Value.(*eth.GetBlockByHashResponse); ok {
			return block, true
		}
	case "earliest":
		return nil, false
	case "pending":
		return nil, false
	default: // hex number
		if !strings.HasPrefix(blockNumberStr, "0x") {
			return nil, false
		}

		blockItr := p.blocks.Back()
		if blockItr == nil {
			return nil, false
		}

		for blockItr != nil {
			if block, ok := blockItr.Value.(*eth.GetBlockByHashResponse); ok {
				requestBlockNumber, err := strconv.ParseInt(blockNumberStr[2:], 16, 64)
				if err != nil {
					return nil, false
				}

				currentBlockNumber, err := strconv.ParseInt(block.Number[2:], 16, 64)
				if err != nil {
					return nil, false
				}

				if requestBlockNumber > currentBlockNumber {
					// request block is higher than last block
					return nil, true
				}

				if currentBlockNumber == requestBlockNumber {
					return block, true
				}
			}

			blockItr = blockItr.Prev()
		}

		// Not found in cache
		return nil, false
	}

	return nil, false
}

func NewBlockPoller(client *qtum.Qtum) (*BlockPoller, error) {
	p := &BlockPoller{client, sync.RWMutex{}, list.New().Init(), false, 256}
	go p.loopSync()
	return p, nil
}
