package transformer

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"

	"github.com/stretchr/testify/assert"
)

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return "0x" + hex.EncodeToString(bytes)
}

type MockBlockPoller struct {
	blocks []*eth.GetBlockByHashResponse
}

func (p *MockBlockPoller) Pull(block *big.Int) (*eth.GetBlockByHashResponse, error) {
	for _, b := range p.blocks {
		bNumber := new(big.Int)
		bNumber.SetString(b.Number[2:], 16)
		if bNumber.Cmp(block) == 0 {
			return b, nil
		}
	}

	return nil, errors.New("Not found block")
}

func (p *MockBlockPoller) addBlock() *eth.GetBlockByHashResponse {
	if p.blocks == nil || len(p.blocks) == 0 {
		p.blocks = []*eth.GetBlockByHashResponse{
			{
				Number:     "0x0",
				Hash:       randomHex(32),
				ParentHash: "0x61d1d322da50dc9961e99d57a14189138065b49046167834bc58e0de84671e14",
			}}

		return p.blocks[0]
	}

	latest := p.blocks[len(p.blocks)-1]
	number := new(big.Int)
	number.SetString(latest.Number[2:], 16)
	number = number.Add(number, big.NewInt(1))

	block := &eth.GetBlockByHashResponse{
		Number:     "0x" + number.Text(16),
		Hash:       randomHex(64),
		ParentHash: latest.Hash,
	}

	p.blocks = append(p.blocks, block)
	return block
}

func initializeBlockPollerAndClient() (*BlockSyncer, internal.Doer, *MockBlockPoller) {
	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, _ := internal.CreateMockedClient(mockedClientDoer)
	poller := &MockBlockPoller{}

	syncer, _ := NewBlockSyncerWithBlockPollerAndInterval(qtumClient, poller, 100*time.Millisecond)
	return syncer, mockedClientDoer, poller
}

func setBlock(doer internal.Doer, poller *MockBlockPoller, start int, end int) {
	if start > len(poller.blocks) {
		panic("Block is not connected")
	}

	poller.blocks = poller.blocks[:start]
	for ; start != end; start++ {
		block := poller.addBlock()
		blockInt := new(big.Int)
		blockInt.SetString(block.Number[2:], 16)

		params, _ := json.Marshal(&qtum.GetBlockHashRequest{
			Int: blockInt,
		})
		doer.AddResponseWithParams(qtum.MethodGetBlockHash, params, qtum.GetBlockHashResponse(block.Hash[2:]))
		doer.AddResponseWithParams(qtum.MethodGetBlockCount, []byte("null"), qtum.GetBlockCountResponse{Int: blockInt})
	}
}

func TestBlockPollerNormalSync(t *testing.T) {
	syncer, doer, poller := initializeBlockPollerAndClient()
	setBlock(doer, poller, 0, 10)

	syncer.Start()
	time.Sleep(1 * time.Millisecond)

	latestBlock, _ := json.Marshal("latest")
	latest, found := syncer.GetBlock(latestBlock)
	latestBlockHash := poller.blocks[len(poller.blocks)-1].Hash

	assert.True(t, found)
	assert.Equal(t, latestBlockHash, latest.Hash)

	blockNineNumber, _ := json.Marshal("0x9")
	blockNine, found := syncer.GetBlock(blockNineNumber)
	blockNineHash := poller.blocks[9].Hash

	assert.True(t, found)
	assert.Equal(t, blockNineHash, blockNine.Hash)
}

func TestBlockPollerSyncUp(t *testing.T) {
	syncer, doer, poller := initializeBlockPollerAndClient()
	setBlock(doer, poller, 0, 10)

	syncer.Start()
	time.Sleep(1 * time.Millisecond)

	latestBlock, _ := json.Marshal("latest")
	latest, found := syncer.GetBlock(latestBlock)
	latestBlockHash := poller.blocks[len(poller.blocks)-1].Hash

	assert.True(t, found)
	assert.Equal(t, latestBlockHash, latest.Hash)

	setBlock(doer, poller, 10, 16)
	time.Sleep(120 * time.Millisecond)

	latestBlock, _ = json.Marshal("latest")
	latest, found = syncer.GetBlock(latestBlock)
	latestBlockHash = poller.blocks[len(poller.blocks)-1].Hash

	assert.True(t, found)
	assert.Equal(t, latestBlockHash, latest.Hash)
	assert.Equal(t, "0xf", latest.Number)
}

func TestBlockPollerForkWithSameHigh(t *testing.T) {
	syncer, doer, poller := initializeBlockPollerAndClient()
	setBlock(doer, poller, 0, 10)

	syncer.Start()
	time.Sleep(1 * time.Millisecond)

	latestBlock, _ := json.Marshal("latest")
	latest, found := syncer.GetBlock(latestBlock)
	latestBlockHash := poller.blocks[len(poller.blocks)-1].Hash

	assert.True(t, found)
	assert.Equal(t, latestBlockHash, latest.Hash)

	setBlock(doer, poller, 8, 10)
	time.Sleep(120 * time.Millisecond)

	latestBlock, _ = json.Marshal("latest")
	newLatest, found := syncer.GetBlock(latestBlock)
	newLatestBlockHash := poller.blocks[len(poller.blocks)-1].Hash

	assert.True(t, found)
	assert.Equal(t, newLatestBlockHash, newLatest.Hash)

	assert.NotEqual(t, latestBlockHash, newLatestBlockHash)
	assert.Equal(t, latest.Number, newLatest.Number)
}

func TestBlockPollerForkWithLowerHight(t *testing.T) {
	syncer, doer, poller := initializeBlockPollerAndClient()
	setBlock(doer, poller, 0, 10)

	syncer.Start()
	time.Sleep(1 * time.Millisecond)

	setBlock(doer, poller, 8, 9)
	time.Sleep(120 * time.Millisecond)

	latestBlock, _ := json.Marshal("latest")
	latest, found := syncer.GetBlock(latestBlock)
	latestBlockHash := poller.blocks[len(poller.blocks)-1].Hash

	assert.True(t, found)
	assert.Equal(t, latestBlockHash, latest.Hash)
	assert.Equal(t, "0x8", latest.Number)

	prevBlock := new(big.Int)
	prevBlock.SetString(latest.Number[2:], 16)
	prevBlock = prevBlock.Sub(prevBlock, big.NewInt(1))

	prevBlockStr := "0x" + prevBlock.Text(16)
	prevBlockJson, _ := json.Marshal(prevBlockStr)
	prev, found := syncer.GetBlock(prevBlockJson)

	assert.False(t, found)
	assert.Nil(t, prev)
}

func TestBlockPollerForkWithHigherHight(t *testing.T) {
	syncer, doer, poller := initializeBlockPollerAndClient()
	setBlock(doer, poller, 0, 10)

	syncer.Start()
	time.Sleep(1 * time.Millisecond)

	setBlock(doer, poller, 8, 16)
	time.Sleep(120 * time.Millisecond)

	latestBlock, _ := json.Marshal("latest")
	latest, found := syncer.GetBlock(latestBlock)
	latestBlockHash := poller.blocks[len(poller.blocks)-1].Hash

	assert.True(t, found)
	assert.Equal(t, latestBlockHash, latest.Hash)
	assert.Equal(t, "0xf", latest.Number)

	prevBlock := new(big.Int)
	prevBlock.SetString(latest.Number[2:], 16)
	prevBlock = prevBlock.Sub(prevBlock, big.NewInt(1))

	prevBlockStr := "0x" + prevBlock.Text(16)
	prevBlockJson, _ := json.Marshal(prevBlockStr)
	prev, found := syncer.GetBlock(prevBlockJson)

	assert.False(t, found)
	assert.Nil(t, prev)
}
