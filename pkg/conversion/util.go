package conversion

import (
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/qtumproject/janus/pkg/utils"
)

func ExtractETHLogsFromTransactionReceipt(blockData qtum.LogBlockData, logs []qtum.Log) []eth.Log {
	result := make([]eth.Log, 0, len(logs))
	for i, log := range logs {
		topics := make([]string, 0, len(log.GetTopics()))
		for _, topic := range log.GetTopics() {
			topics = append(topics, utils.AddHexPrefix(topic))
		}
		result = append(result, eth.Log{
			TransactionHash:  utils.AddHexPrefix(blockData.GetTransactionHash()),
			TransactionIndex: hexutil.EncodeUint64(blockData.GetTransactionIndex()),
			BlockHash:        utils.AddHexPrefix(blockData.GetBlockHash()),
			BlockNumber:      hexutil.EncodeUint64(blockData.GetBlockNumber()),
			Data:             utils.AddHexPrefix(log.GetData()),
			Address:          utils.AddHexPrefix(log.GetAddress()),
			Topics:           topics,
			LogIndex:         hexutil.EncodeUint64(uint64(i)),
		})
	}
	return result
}
