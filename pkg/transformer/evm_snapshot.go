package transformer

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

// Implements https://github.com/trufflesuite/ganache-cli#custom-methods

// chainSnapshot contains information to emulate evm_snapshot provided by ganache
type chainSnapshot struct {
	Height    int
	BlockHash string
}

func (s *chainSnapshot) ID() string {
	// let's just height... simpler to manipulate by hand, and we don't expect
	// reorgs in regtest
	return fmt.Sprintf("%d", s.Height)
	// return fmt.Sprintf("%d:%s", s.Height, s.BlockHash)
}

// janusChainSnapshots is a singletone used to store snapshots taken with evm_snapshot
var janusChainSnapshots chainSnapshots

type chainSnapshots struct {
	sync.Mutex

	items []chainSnapshot
}

func (s *chainSnapshots) RevertTo(id string) (chainSnapshot, bool) {
	s.Lock()
	defer s.Unlock()

	// Revert the state of the blockchain to a previous snapshot. Takes a single
	// parameter, which is the snapshot id to revert to. This deletes the given
	// snapshot, as well as any snapshots taken after.

	for i := len(s.items) - 1; i >= 0; i-- {
		snapshot := s.items[i]

		if snapshot.ID() == id {
			s.items = s.items[:i]
			return snapshot, true
		}
	}

	return chainSnapshot{}, false
}

func (s *chainSnapshots) Add(newSnapshot chainSnapshot) {
	s.Lock()
	defer s.Unlock()

	if len(s.items) > 0 {
		latestSnapshot := s.items[len(s.items)-1]
		if latestSnapshot == newSnapshot {
			return
		}
	}

	s.items = append(s.items, newSnapshot)
}

type ProxyEVMSnapshot struct {
	*qtum.Qtum
}

func (p *ProxyEVMSnapshot) Method() string {
	return "evm_snapshot"
}

func (p *ProxyEVMSnapshot) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	// Mine a block before taking the snapshot. We want to restore to the block
	// before this generated block. When revert is called on this snapshot id,
	// this block will be invalidated.
	//
	// Note: there might be a corner case where this generated block is not empty.
	// To avoid this problem, should only use janus in regtest by itself, s.t.
	// every tx is mined immediately with Generate(1).
	_, err := p.Generate(1)
	if err != nil {
		return nil, err
	}

	var qinfo qtum.GetBlockChainInfoResponse
	if err := p.Qtum.Request(qtum.MethodGetBlockChainInfo, nil, &qinfo); err != nil {
		return nil, err
	}

	snapshot := chainSnapshot{
		Height:    int(qinfo.Blocks),
		BlockHash: qinfo.Bestblockhash,
	}

	janusChainSnapshots.Add(snapshot)

	spew.Dump(janusChainSnapshots.items)

	return eth.EVMSnapshotResponse(snapshot.ID()), nil
}

type ProxyEVMRevert struct {
	*qtum.Qtum
}

func (p *ProxyEVMRevert) Method() string {
	return "evm_revert"
}

func (p *ProxyEVMRevert) Request(rpcreq *eth.JSONRPCRequest) (interface{}, error) {
	var snapshotID eth.EVMRevertRequest
	if err := json.Unmarshal(rpcreq.Params, &snapshotID); err != nil {
		return nil, err
	}

	snapshot, ok := janusChainSnapshots.RevertTo(string(snapshotID))

	if !ok {
		return nil, errors.New("snapshot not found")
	}

	_, err := p.Qtum.InvalidateBlock(snapshot.BlockHash)
	if err != nil {
		return nil, err
	}

	log.Println("reverted", snapshot)
	spew.Dump(janusChainSnapshots.items)

	return eth.EVMRevertResponse(true), nil
}
