package enums

type PairEvent string

const (
	PairEventSwap PairEvent = "Swap"
	PairEventSync PairEvent = "Sync"
	PairEventMint PairEvent = "Mint"
	PairEventBurn PairEvent = "Burn"
)

func (e PairEvent) ToString() string {
	return string(e)
}

type SwapTxSide string

const (
	SwapBuy  SwapTxSide = "BUY"
	SwapSell SwapTxSide = "SELL"
)
