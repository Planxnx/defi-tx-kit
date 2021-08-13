package utils

import (
	"log"
	"strings"

	"github.com/Planxnx/defi-tx-kit/contract"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	PairABI  abi.ABI = MustParseABI(contract.PairABI)
	TokenABI abi.ABI = MustParseABI(contract.TokenABI)
)

func MustParseABI(strABI string) abi.ABI {
	abi, err := abi.JSON(strings.NewReader(strABI))
	if err != nil {
		log.Fatal("Can't parse contract abi ", err)
	}
	return abi
}
