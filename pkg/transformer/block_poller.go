package transformer

import (
	"container/list"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

type BlockPoller interface {
	Pull(blockNumber *big.Int) (*eth.GetBlockByHashResponse, error)
}

type DefaultBlockPoller struct {
	Qtum *qtum.Qtum
}

func (p *DefaultBlockPoller) Pull(blockNumber *big.Int) (*eth.GetBlockByHashResponse, error) {
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

type BlockSyncer struct {
	Qtum     *qtum.Qtum
	lock     sync.RWMutex
	blocks   *list.List
	synced   bool
	limit    int
	poller   BlockPoller
	interval time.Duration
}

func (s *BlockSyncer) clearBlocks() {
	s.synced = false

	// Assume lock is locked
	itr := s.blocks.Front()

	// No valid block then cleanup
	for itr != nil {
		delItr := itr
		itr = itr.Next()
		s.blocks.Remove(delItr)
	}
}

func (s *BlockSyncer) loopSync() error {

	for {
		// Query block count
		blockCountResp, err := s.Qtum.GetBlockCount()
		if err != nil {
			s.Qtum.GetLogger().Log("function", "loopSync", "message", "fail to query blockcount", "error", err)
			continue
		}
		upstreamBlock := blockCountResp.Int

		// Get local block
		blockItr := s.blocks.Back()

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
			newBlock, err := s.poller.Pull(big.NewInt(0).Add(localBlock, big.NewInt(1)))
			if err != nil {
				s.Qtum.GetLogger().Log("function", "loopSync", "message", "fail to query block", "error", err)
				continue
			}

			if localHash == "" || newBlock.ParentHash == localHash {
				s.lock.Lock()
				s.blocks.PushBack(newBlock)
				s.synced = true
				s.lock.Unlock()
			} else {
				s.Qtum.GetLogger().Log("function", "loopSync", "message", "last block is invalid", "error", err)
				s.lock.Lock()
				s.clearBlocks()
				s.lock.Unlock()
				continue
			}
		} else if localBlock.Cmp(upstreamBlock) > 0 {
			s.lock.Lock()
			s.clearBlocks()
			s.lock.Unlock()
			continue
		} else {
			upstreamHash, err := s.Qtum.GetBlockHash(localBlock)
			if err != nil {
				s.Qtum.GetLogger().Log("function", "loopSync", "message", "Fail to get block hash of local height", "error", err)
				continue
			}

			if "0x"+upstreamHash != qtum.GetBlockHashResponse(localHash) {
				s.Qtum.GetLogger().Log("function", "loopSync", "message", "Upstream hash is not match with local", "error", err)
				s.lock.Lock()
				s.clearBlocks()
				s.lock.Unlock()
				continue
			}

			// If last block is corrent just sleep
			time.Sleep(s.interval)
		}

		// Cleanup old
		for s.blocks.Len() > s.limit {
			s.lock.Lock()
			s.blocks.Remove(s.blocks.Front())
			s.lock.Unlock()
		}
	}
}

func (s *BlockSyncer) GetBlock(blockNumber json.RawMessage) (*eth.GetBlockByHashResponse, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.synced {
		return nil, false
	}

	var blockNumberStr string
	if err := json.Unmarshal(blockNumber, &blockNumberStr); err != nil {
		return nil, false
	}

	var element *list.Element
	switch blockNumberStr {
	case "latest":
		element = s.blocks.Back()
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

		blockItr := s.blocks.Back()
		if blockItr == nil {
			return nil, false
		}

		for blockItr != nil {
			if block, ok := blockItr.Value.(*eth.GetBlockByHashResponse); ok {
				requestBlockNumber := new(big.Int)
				if _, ok := requestBlockNumber.SetString(blockNumberStr[2:], 16); !ok {
					fmt.Println("Fail to decode", blockNumberStr[2:])
					return nil, false
				}

				currentBlockNumber := new(big.Int)
				if _, ok := currentBlockNumber.SetString(block.Number[2:], 16); !ok {
					fmt.Println("Fail to decode", block.Number[2:])
					return nil, false
				}

				fmt.Println(requestBlockNumber, currentBlockNumber)

				if requestBlockNumber.Cmp(currentBlockNumber) > 0 {
					// request block is higher than last block
					return nil, true
				}

				if currentBlockNumber.Cmp(requestBlockNumber) == 0 {
					return block, true
				}
			}

			blockItr = blockItr.Prev()
		}
	}

	return nil, false
}

func (s *BlockSyncer) Start() {
	go s.loopSync()
}

func NewBlockSyncer(client *qtum.Qtum) (*BlockSyncer, error) {
	return NewBlockSyncerWithBlockPollerAndInterval(client, &DefaultBlockPoller{client}, 200*time.Millisecond)
}

func NewBlockSyncerWithBlockPollerAndInterval(client *qtum.Qtum, poller BlockPoller, interval time.Duration) (*BlockSyncer, error) {
	s := &BlockSyncer{client, sync.RWMutex{}, list.New().Init(), false, 256, poller, interval}
	return s, nil
}
