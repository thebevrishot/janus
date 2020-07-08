package transformer

import (
	"encoding/json"
	"fmt"
	"math/big"

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
	gasLimit = big.NewInt(40000000)
	if gas := g.GasHex(); gas != "" {
		gasLimit, err = utils.DecodeBig(gas)
		if err != nil {
			err = errors.Wrap(err, "decode gas")
			return
		}
	}

	gasPriceFloat64, err := EthValueToQtumAmount(g.GasPriceHex())
	if err != nil {
		return nil, "0.0", err
	}
	gasPrice = fmt.Sprintf("%.8f", gasPriceFloat64)

	return
}

func EthValueToQtumAmount(val string) (float64, error) {
	if val == "" {
		return 0.0000004, nil
	}

	ethVal, err := utils.DecodeBig(val)
	if err != nil {
		return 0.0, err
	}

	ethDecimal := decimal.NewFromBigInt(ethVal, 0)

	amount := ethDecimal.Mul(decimal.NewFromFloat(float64(1e-8)))
	result, exact := amount.Float64()
	if !exact {
		return 0, fmt.Errorf("Could not get an exact value for value %v got %v", val, result)
	}
	return result, nil
}

func QtumAmountToEthValue(amount float64) (string, error) {
	bigAmount := decimal.NewFromFloat(amount)
	bigAmount = bigAmount.Mul(decimal.NewFromFloat(1e8))

	result := new(big.Int)
	result, success := result.SetString(bigAmount.String(), 10)
	if !success {
		return "0x0", errors.New("big.Int#SetString is not success")
	}

	return hexutil.EncodeBig(result), nil
}

func unmarshalRequest(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(err, "Invalid RPC input")
	}
	return nil
}
