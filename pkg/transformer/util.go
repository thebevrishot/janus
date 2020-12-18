package transformer

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/qtumproject/janus/pkg/eth"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/utils"
	"github.com/shopspring/decimal"
)

type EthGas interface {
	GasHex() string
	GasPriceHex() string
}

func EthGasToQtum(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = g.(*eth.SendTransactionRequest).Gas.Int

	gasPriceDecimal, err := EthValueToQtumAmount(g.GasPriceHex())
	if err != nil {
		return nil, "0.0", err
	}
	gasPrice = fmt.Sprintf("%v", gasPriceDecimal)

	return
}

func EthValueToQtumAmount(val string) (decimal.Decimal, error) {
	if val == "" {
		return decimal.NewFromFloat(0.0000004), nil
	}

	ethVal, err := utils.DecodeBig(val)
	if err != nil {
		return decimal.NewFromFloat(0.0), err
	}

	ethValDecimal, err := decimal.NewFromString(ethVal.String())
	if err != nil {
		return decimal.NewFromFloat(0.0), errors.New("decimal.NewFromString was not a success")
	}

	amount := ethValDecimal.Mul(decimal.NewFromFloat(float64(1e-8)))

	return amount, nil
}

func formatQtumAmount(amount decimal.Decimal) (string, error) {
	decimalAmount := amount.Mul(decimal.NewFromFloat(float64(1e8)))

	//convert decimal to Integer
	result := decimalAmount.BigInt()

	if !decimalAmount.Equals(decimal.NewFromBigInt(result, 0)) {
		return "0x0", errors.New("decimal.BigInt() was not a success")
	}

	return hexutil.EncodeBig(result), nil
}

func unmarshalRequest(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(err, "Invalid RPC input")
	}
	return nil
}
