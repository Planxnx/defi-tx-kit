package utils

import (
	"github.com/Planxnx/defi-tx-kit/contract"
	"github.com/Planxnx/defi-tx-kit/enums"
	"github.com/ethereum/go-ethereum/core/types"
)

func PairParseSync(txLog types.Log) (*contract.PairSync, error) {
	event := &contract.PairSync{}
	if err := PairABI.UnpackIntoInterface(event, enums.PairEventSync.ToString(), txLog.Data); err != nil {
		return nil, err
	}
	event.Raw = txLog
	return event, nil
}

func PairParseSwap(txLog types.Log) (*contract.PairSwap, error) {
	event := &contract.PairSwap{}
	if err := PairABI.UnpackIntoInterface(event, enums.PairEventSwap.ToString(), txLog.Data); err != nil {
		return nil, err
	}
	event.Raw = txLog
	return event, nil
}
