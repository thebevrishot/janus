package transformer

import (
	"fmt"
	"testing"

	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestEthValueToQtumAmount(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"in":   "0x64",
			"want": decimal.NewFromFloat(0.000001),
		},
		{

			"in":   "0x1",
			"want": decimal.NewFromFloat(0.00000001),
		},
	}
	for _, c := range cases {
		in := c["in"].(string)
		want := c["want"].(decimal.Decimal)
		got, err := EthValueToQtumAmount(in)
		if err != nil {
			t.Error(err)
		}
		if !got.Equal(want) {
			t.Errorf("in: %s, want: %v, got: %v", in, want, got)
		}
	}
}

func TestQtumAmountToEthValue(t *testing.T) {
	in, want := decimal.NewFromFloat(0.000001), "0x64"
	got, err := formatQtumAmount(in)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("in: %v, want: %s, got: %s", in, want, got)
	}
}

func TestAddressesConvertion(t *testing.T) {
	t.Parallel()

	inputs := []struct {
		qtumChain   string
		ethAddress  string
		qtumAddress string
	}{
		{
			qtumChain:   qtum.ChainTest,
			ethAddress:  "6c89a1a6ca2ae7c00b248bb2832d6f480f27da68",
			qtumAddress: "qTTH1Yr2eKCuDLqfxUyBLCAjmomQ8pyrBt",
		},
	}

	for i, in := range inputs {
		var (
			in       = in
			testDesc = fmt.Sprintf("#%d", i)
		)
		t.Run(testDesc, func(t *testing.T) {
			qtumAddress, err := convertETHAddress(in.ethAddress, in.qtumChain)
			require.NoError(t, err, "couldn't convert Ethereum address to Qtum address")
			require.Equal(t, in.qtumAddress, qtumAddress, "unexpected converted Qtum address value")

			ethAddress, err := convertQtumAddress(in.qtumAddress)
			require.NoError(t, err, "couldn't convert Qtum address to Ethereum address")
			require.Equal(t, in.ethAddress, ethAddress, "unexpected converted Ethereum address value")
		})
	}
}
