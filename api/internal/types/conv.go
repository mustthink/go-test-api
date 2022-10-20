package types

import (
	"math"
	"math/big"
	"strconv"
	"strings"
)

func ConvHexDec(s string) (uint64, error) {
	s = strings.Replace(s, "0x", "", -1)
	dec, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, err
	}
	return dec, nil
}

func ConvHexWeiToDecEth(s string) float64 {
	dec := new(big.Int)
	s = strings.Replace(s, "0x", "", -1)
	dec.SetString(s, 16)

	bigval := new(big.Float).SetInt(dec)
	bigval = new(big.Float).Quo(bigval, big.NewFloat(math.Pow(10, 18)))

	value, _ := bigval.Float64()
	return value
}

func ConvDecHex(d int64) string {
	s := strconv.FormatInt(d, 16)
	return "0x" + s
}
