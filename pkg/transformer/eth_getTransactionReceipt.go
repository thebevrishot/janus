package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/conversion"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"
)

// ProxyETHGetTransactionReceipt implements ETHProxy
type ProxyETHGetTransactionReceipt struct {
	*qtum.Qtum
}

func (p *ProxyETHGetTransactionReceipt) Method() string {
	return "eth_getTransactionReceipt"
}

func (p *ProxyETHGetTransactionReceipt) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, error) {
	var req eth.GetTransactionReceiptRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}
	if req == "" {
		return nil, errors.New("empty transaction hash")
	}
	var (
		txHash  = utils.RemoveHexPrefix(string(req))
		qtumReq = qtum.GetTransactionReceiptRequest(txHash)
	)
	return p.request(&qtumReq)
}

func (p *ProxyETHGetTransactionReceipt) request(req *qtum.GetTransactionReceiptRequest) (*eth.GetTransactionReceiptResponse, error) {
	isVMTransaction := false
	qtumReceipt, err := p.Qtum.GetTransactionReceipt(string(*req))
	if err == nil {
		isVMTransaction = true
	} else {
		if qtumReceipt == nil {
			// panic(err)
		}
		rawTransaction, getRawTransactionErr := p.Qtum.GetRawTransaction(string(*req), false)
		if getRawTransactionErr == nil {
			block, err := p.Qtum.GetBlock(rawTransaction.BlockHash)
			if err != nil {
				p.Qtum.GetDebugLogger().Log("msg", "Failed to get block with hash", "hash", rawTransaction.BlockHash, "err", err)
				return nil, err
			}
			p.Qtum.GetDebugLogger().Log("msg", "Transaction does not execute VM code so does not have a transaction receipt, returning a dummy", "txid", string(*req))
			qtumReceipt = &qtum.GetTransactionReceiptResponse{
				TransactionHash:  string(*req),
				TransactionIndex: 1,
				BlockHash:        rawTransaction.BlockHash,
				BlockNumber:      uint64(block.Height),
			}
		} else {
			p.Qtum.GetDebugLogger().Log("msg", "Transaction does not exist", "txid", string(*req))
			errCause := errors.Cause(err)
			if errCause == qtum.EmptyResponseErr {
				return nil, nil
			}
			return nil, err
		}
	}

	ethReceipt := &eth.GetTransactionReceiptResponse{
		TransactionHash:   utils.AddHexPrefix(qtumReceipt.TransactionHash),
		TransactionIndex:  hexutil.EncodeUint64(qtumReceipt.TransactionIndex),
		BlockHash:         utils.AddHexPrefix(qtumReceipt.BlockHash),
		BlockNumber:       hexutil.EncodeUint64(qtumReceipt.BlockNumber),
		ContractAddress:   utils.AddHexPrefixIfNotEmpty(qtumReceipt.ContractAddress),
		CumulativeGasUsed: hexutil.EncodeUint64(qtumReceipt.CumulativeGasUsed),
		GasUsed:           hexutil.EncodeUint64(qtumReceipt.GasUsed),
		From:              utils.AddHexPrefixIfNotEmpty(qtumReceipt.From),
		To:                utils.AddHexPrefixIfNotEmpty(qtumReceipt.To),

		// TODO: researching
		// ! Temporary accept this value to be always zero, as it is at eth logs
		LogsBloom: eth.EmptyLogsBloom,
	}

	status := "0x0"
	if qtumReceipt.Excepted == "None" {
		status = "0x1"
	}
	ethReceipt.Status = status

	if isVMTransaction {
		r := qtum.TransactionReceipt(*qtumReceipt)
		ethReceipt.Logs = conversion.ExtractETHLogsFromTransactionReceipt(&r)

		qtumTx, err := p.Qtum.GetTransaction(qtumReceipt.TransactionHash)
		if err != nil {
			return nil, errors.WithMessage(err, "couldn't get transaction")
		}
		decodedRawQtumTx, err := p.Qtum.DecodeRawTransaction(qtumTx.Hex)
		if err != nil {
			return nil, errors.WithMessage(err, "couldn't decode raw transaction")
		}
		if decodedRawQtumTx.IsContractCreation() {
			ethReceipt.To = ""
		} else {
			ethReceipt.ContractAddress = ""
		}
	}

	// TODO: researching
	// - The following code reason is unknown (see original comment)
	// - Code temporary commented, until an error occures
	// ! Do not remove
	// // contractAddress : DATA, 20 Bytes - The contract address created, if the transaction was a contract creation, otherwise null.
	// if status != "0x1" {
	// 	// if failure, should return null for contractAddress, instead of the zero address.
	// 	ethTxReceipt.ContractAddress = ""
	// }

	return ethReceipt, nil
}
