package utils

import (
	"fmt"
	"math"
	"math/big"

	"github.com/Planxnx/defi-tx-kit/contract"
	"github.com/Planxnx/defi-tx-kit/enums"
	"github.com/ethereum/go-ethereum/core/types"
)

func TokenParseTransfer(txLog types.Log) (*contract.TokenTransfer, error) {
	event := &contract.TokenTransfer{}
	if err := TokenABI.UnpackIntoInterface(event, enums.TokenEventTransfer.ToString(), txLog.Data); err != nil {
		return nil, err
	}
	event.Raw = txLog
	return event, nil
}

func TokenParseApproval(txLog types.Log) (*contract.TokenApproval, error) {
	event := &contract.TokenApproval{}
	if err := TokenABI.UnpackIntoInterface(event, enums.TokenEventApproval.ToString(), txLog.Data); err != nil {
		return nil, err
	}
	event.Raw = txLog
	return event, nil
}

func ConvertToDecimals(amount *big.Int, decimals int) *big.Float {
	decimalsFloat := big.NewFloat(math.Pow10(decimals))
	amountFloat := big.NewFloat(0).SetInt(amount)
	return big.NewFloat(0).Quo(amountFloat, decimalsFloat)
}

func ConvertFromDecimals(amount *big.Float, decimals int) *big.Int {
	decimalsFloat := big.NewFloat(math.Pow10(decimals))
	value := big.NewFloat(0).Mul(amount, decimalsFloat)
	valueInt64, _ := value.Int64()
	if valueInt64 == math.MaxInt64 || valueInt64 == math.MinInt64 {
		panic(fmt.Sprintf("convertFromDecimals: amount.(*big.Float) = %v shoud more than math.MinInt64 or less thane math.MaxInt64\n", amount))
	}
	return big.NewInt(valueInt64)
}
