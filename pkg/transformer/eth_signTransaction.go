package transformer

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/utils"

	"github.com/shopspring/decimal"
)

// ProxyETHSendTransaction implements ETHProxy
type ProxyETHSignTransaction struct {
	*qtum.Qtum
}

func (p *ProxyETHSignTransaction) Method() string {
	return "eth_signTransaction"
}

func (p *ProxyETHSignTransaction) Request(rawreq *eth.JSONRPCRequest) (interface{}, error) {
	var req eth.SendTransactionRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}

	fixedAmount := decimal.NewFromBigInt(req.Value.Int, 0)

	// get necessary utxo ids needed for creating raw transaction
	inputs, err := p.getRequiredUtxos(req.From, fixedAmount)
	if err != nil {
		return nil, err
	}

	if req.IsCreateContract() {
		return p.requestCreateContract(&req, inputs)
	} else if req.IsSendEther() {
		return p.requestSendToAddress(&req, inputs)
	} else if req.IsCallContract() {
		return p.requestSendToContract(&req, inputs)
	}

	return nil, errors.New("Unknown operation")
}

func (p *ProxyETHSignTransaction) getRequiredUtxos(from string, neededAmount decimal.Decimal) ([]qtum.RawTxInputs, error) {
	//convert address to qtum address
	addr := utils.RemoveHexPrefix(from)
	base58Addr, err := p.FromHexAddress(addr)
	if err != nil {
		return nil, err
	}
	// need to get utxos with txid and vouts. In order to do this we get a list of unspent transactions and begin summing them up
	// todo: convert from to actual qtum address/figure out how to do that.
	var unspentListReq *qtum.ListUnspentRequest = &qtum.ListUnspentRequest{6, 999, []string{base58Addr}}
	qtumresp, err := p.ListUnspent(unspentListReq)
	if err != nil {
		return nil, err
	}

	balance := decimal.New(0, 0)
	var inputs []qtum.RawTxInputs
	var balanceReqMet bool
	for _, utxo := range *qtumresp {
		balance.Add(decimal.NewFromFloat(utxo.Amount))
		inputs = append(inputs, qtum.RawTxInputs{utxo.Txid, utxo.Vout})
		if balance.GreaterThanOrEqual(neededAmount) {
			balanceReqMet = true
			break
		}
	}
	if balanceReqMet {
		return inputs, nil
	}
	return nil, fmt.Errorf("Insufficient UTXO value attempted to be sent")
}

func (p *ProxyETHSignTransaction) requestSendToContract(ethtx *eth.SendTransactionRequest, inputs []qtum.RawTxInputs) (string, error) {
	gasLimit, gasPrice, err := EthGasToQtum(ethtx)
	if err != nil {
		return "", err
	}

	var amount decimal.Decimal
	if ethtx.Value.Cmp(big.NewInt(0)) != 0 {
		var err error
		amount, err = EthValueToQtumAmount(ethtx.Value.Int)
		if err != nil {
			return "", errors.Wrap(err, "EthValueToQtumAmount:")
		}
	}

	contractInteractTx := qtum.SendToContractRequest{
		ContractAddress: utils.RemoveHexPrefix(ethtx.To),
		Datahex:         utils.RemoveHexPrefix(ethtx.Data),
		Amount:          amount.StringFixed(8),
		GasLimit:        gasLimit,
		GasPrice:        gasPrice,
	}

	if from := ethtx.From; from != "" && utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return "", err
		}
		contractInteractTx.SenderAddress = from
	}

	fromAddr := utils.RemoveHexPrefix(ethtx.From)

	acc := p.Qtum.Accounts.FindByHexAddress(fromAddr)
	if acc == nil {
		return "", errors.Errorf("No such account: %s", fromAddr)
	}

	var contractMap = map[string]interface{}{"contract": contractInteractTx}
	rawtxreq := []interface{}{inputs, contractMap}
	var rawTx string
	if err := p.Qtum.Request(qtum.MethodCreateRawTx, rawtxreq, &rawTx); err != nil {
		return "", err
	}

	var resp *qtum.SignRawTxResponse
	var privkeyArray = []string{acc.String()}
	signrawtxreq := []interface{}{rawTx, privkeyArray}
	if err := p.Qtum.Request(qtum.MethodSignRawTx, signrawtxreq, &resp); err != nil {
		return "", err
	}
	if len(resp.Errors) != 0 {
		var errStr = []string{"List of errors in raw transaction signing: "}
		for _, i := range resp.Errors {
			errStr = append(errStr, fmt.Sprint("For txid %v there was an error: %v", i.Txid, i.Error))
		}
		return "", fmt.Errorf(strings.Join(errStr, "\n"))
	}
	return resp.Hex, nil
}

func (p *ProxyETHSignTransaction) requestSendToAddress(req *eth.SendTransactionRequest, inputs []qtum.RawTxInputs) (string, error) {
	getQtumWalletAddress := func(addr string) (string, error) {
		if utils.IsEthHexAddress(addr) {
			return p.FromHexAddress(utils.RemoveHexPrefix(addr))
		}
		return addr, nil
	}

	to, err := getQtumWalletAddress(req.To)
	if err != nil {
		return "", err
	}

	fromAddr := utils.RemoveHexPrefix(req.From)

	acc := p.Qtum.Accounts.FindByHexAddress(fromAddr)
	if acc == nil {
		return "", errors.Errorf("No such account: %s", fromAddr)
	}

	amount, err := EthValueToQtumAmount(req.Value.Int)
	if err != nil {
		return "", err
	}

	var addressValMap = map[string]interface{}{to: amount}
	rawtxreq := []interface{}{inputs, addressValMap}
	var rawTx string
	if err := p.Qtum.Request(qtum.MethodCreateRawTx, rawtxreq, &rawTx); err != nil {
		return "", err
	}

	var resp *qtum.SignRawTxResponse
	var privkeyArray = []string{acc.String()}
	signrawtxreq := []interface{}{rawTx, privkeyArray}
	if err := p.Qtum.Request(qtum.MethodSignRawTx, signrawtxreq, &resp); err != nil {
		return "", err
	}
	if len(resp.Errors) != 0 {
		var errStr = []string{"List of errors in raw transaction signing: "}
		for _, i := range resp.Errors {
			errStr = append(errStr, fmt.Sprint("For txid %v there was an error: %v", i.Txid, i.Error))
		}
		return "", fmt.Errorf(strings.Join(errStr, "\n"))
	}
	return resp.Hex, nil
}

func (p *ProxyETHSignTransaction) requestCreateContract(req *eth.SendTransactionRequest, inputs []qtum.RawTxInputs) (string, error) {
	gasLimit, gasPrice, err := EthGasToQtum(req)
	if err != nil {
		return "", err
	}

	contractDeploymentTx := &qtum.CreateContractRequest{
		ByteCode: utils.RemoveHexPrefix(req.Data),
		GasLimit: gasLimit,
		GasPrice: gasPrice,
	}

	from := req.From
	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return "", err
		}
	}

	contractDeploymentTx.SenderAddress = from

	fromAddr := utils.RemoveHexPrefix(req.From)

	acc := p.Qtum.Accounts.FindByHexAddress(fromAddr)
	if acc == nil {
		return "", errors.Errorf("No such account: %s", fromAddr)
	}

	var contractMap = map[string]interface{}{"contract": contractDeploymentTx}
	rawtxreq := []interface{}{inputs, contractMap}
	var rawTx string
	if err := p.Qtum.Request(qtum.MethodCreateRawTx, rawtxreq, &rawTx); err != nil {
		return "", err
	}

	var resp *qtum.SignRawTxResponse
	var privkeyArray = []string{acc.String()}
	signrawtxreq := []interface{}{rawTx, privkeyArray}
	if err := p.Qtum.Request(qtum.MethodSignRawTx, signrawtxreq, &resp); err != nil {
		return "", err
	}
	if len(resp.Errors) != 0 {
		var errStr = []string{"List of errors in raw transaction signing: "}
		for _, i := range resp.Errors {
			errStr = append(errStr, fmt.Sprint("For txid %v there was an error: %v", i.Txid, i.Error))
		}
		return "", fmt.Errorf(strings.Join(errStr, "\n"))
	}
	return resp.Hex, nil
}
