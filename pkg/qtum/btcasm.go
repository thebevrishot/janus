package qtum

import (
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

type (
	// ASM is Bitcoin Script extended by Qtum to support smart contracts
	ASM struct {
		VMVersion   string
		GasLimitStr string
		GasPriceStr string

		// OP_CREATE || OP_CALL
		Instructor string // FIXME: typo
	}
	CallASM struct {
		ASM
		callData        string
		ContractAddress string
	}
	CreateASM struct {
		ASM
		callData string
	}

	CreateSenderASM struct {
		ASM
		Sender   string
		CallData string
	}

	ContractInvokeInfo struct {
		// VMVersion string
		From     string
		GasLimit string
		GasPrice string
		CallData string
	}
)

func (asm *ASM) GasPrice() (*big.Int, error) {
	return stringNumberToBigInt(asm.GasPriceStr)
}

func (asm *ASM) GasLimit() (*big.Int, error) {
	return stringNumberToBigInt(asm.GasLimitStr)
}

func (asm *CreateASM) CallData() string {
	return asm.callData
}

func (asm *CallASM) CallData() string {
	return asm.callData
}

func ParseCreateSenderASM(asm string) (*ContractInvokeInfo, error) {
	// See: https://github.com/qtumproject/qips/issues/6

	// "1 7926223070547d2d15b2ef5e7383e541c338ffe9 6a473044022067ca66b0308ae16aeca7a205ce0490b44a61feebe5632710b52aabde197f9e4802200e8beec61a58dbe1279a9cdb68983080052ae7b9997bc863b7c5623e4cb55fd
	// b01210299d391f528b9edd07284c7e23df8415232a8ce41531cf460a390ce32b4efd112 OP_SENDER 4 6721975 100 6060604052341561000f57600080fd5b60008054600160a060020a033316600160a060020a03199091161790556101de8061003b6000
	// 396000f300606060405263ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630900f010811461005d578063445df0ac1461007e5780638da5cb5b146100a3578063fdacd576146100d257600080fd5b341561
	// 006857600080fd5b61007c600160a060020a03600435166100e8565b005b341561008957600080fd5b61009161017d565b60405190815260200160405180910390f35b34156100ae57600080fd5b6100b6610183565b604051600160a060020a039091168152
	// 60200160405180910390f35b34156100dd57600080fd5b61007c600435610192565b6000805433600160a060020a03908116911614156101795781905080600160a060020a031663fdacd5766001546040517c01000000000000000000000000000000000000
	// 0000000000000000000063ffffffff84160281526004810191909152602401600060405180830381600087803b151561016457600080fd5b6102c65a03f1151561017557600080fd5b5050505b5050565b60015481565b600054600160a060020a031681565b
	// 60005433600160a060020a03908116911614156101af5760018190555b505600a165627a7a72305820b6a912c5b5115d1a5412235282372dc4314f325bac71ee6c8bd18f658d7ed1ad0029 OP_CREATE"

	parts := strings.Split(asm, " ")
	if len(parts) < 9 {
		return nil, errors.New("invalid create_sender script")
	}

	gasLimit, err := stringBase10ToHex(parts[5])
	if err != nil {
		return nil, err
	}

	gasPrice, err := stringBase10ToHex(parts[6])
	if err != nil {
		return nil, err
	}

	return &ContractInvokeInfo{
		From:     parts[1],
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		CallData: parts[7],
	}, nil
}

func ParseCreateASM(asm string) (*CreateASM, error) {
	parts := strings.Split(asm, " ")
	if len(parts) < 5 {
		return nil, errors.New("invalid create ASM")
	}

	return &CreateASM{
		ASM: ASM{
			VMVersion:   parts[0],
			GasLimitStr: parts[1],
			GasPriceStr: parts[2],
			Instructor:  parts[4],
		},
		callData: parts[3],
	}, nil
}

func ParseCallASM(asm string) (*CallASM, error) {
	parts := strings.Split(asm, " ")
	if len(parts) < 6 {
		return nil, errors.New("invalid call ASM")
	}

	return &CallASM{
		ASM: ASM{
			VMVersion:   parts[0],
			GasLimitStr: parts[1],
			GasPriceStr: parts[2],
			Instructor:  parts[5],
		},
		callData:        parts[3],
		ContractAddress: parts[4],
	}, nil
}

func stringBase10ToHex(str string) (string, error) {
	var v big.Int
	_, ok := v.SetString(str, 10)
	if !ok {
		return "", errors.Errorf("Failed to parse big.Int: %s", str)
	}

	return v.Text(16), nil
}

func stringNumberToBigInt(str string) (*big.Int, error) {
	var success bool
	v := new(big.Int)
	if v, success = v.SetString(str, 10); !success {
		return nil, errors.Errorf("Failed to parse big.Int: %s", str)
	}
	return v, nil
}
