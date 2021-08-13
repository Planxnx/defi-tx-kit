package utils

import (
	"context"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Chain string

type ChainData struct {
	// Network Settings
	id        int
	name      string
	currency  string
	jsonRPC   string
	wsRPC     string
	blocktime int

	// RPC Client
	ethclient   *ethclient.Client
	ethwsclient *ethclient.Client
	once        sync.Once
}

const (
	ChainEthereum           Chain = "ethereum"
	ChainOptimisticEthereum Chain = "optimisticethereum"
	ChainBinanceSmartChain  Chain = "bsc"
	ChainOptimisticKovan    Chain = "optimistickovan"
	ChainPolygon            Chain = "polygon"
)

func IsValidChain(c *Chain) bool {
	_, ok := chainDataMap[*c]
	return ok
}

func (c Chain) ID() int {
	return c.data().id
}

func (c Chain) Name() string {
	return c.data().name
}

func (c Chain) Currency() string {
	return c.data().currency
}

func (c Chain) GetBlockTime() int {
	return c.data().blocktime
}

func (c Chain) GetClient(ctx context.Context) (client *ethclient.Client, err error) {
	chainData := c.data()
	if chainData.jsonRPC == "" {
		return nil, errors.New("undefined rpc endpoint")
	}
	chainData.once.Do(func() {
		chainData.ethclient, err = ethclient.DialContext(ctx, chainData.jsonRPC)
	})
	return chainData.ethclient, err
}

func (c Chain) GetWSClient(ctx context.Context) (client *ethclient.Client, err error) {
	chainData := c.data()
	if chainData.wsRPC == "" {
		return nil, errors.New("undefined rpc endpoint")
	}
	chainData.once.Do(func() {
		chainData.ethwsclient, err = ethclient.DialContext(ctx, chainData.wsRPC)
	})
	return chainData.ethwsclient, err
}

func (c *Chain) data() *ChainData {
	data, ok := chainDataMap[*c]
	if !ok {
		panic("invalid chain")
	}
	return data
}

var chainDataMap = map[Chain]*ChainData{
	ChainEthereum:           chainEthereumData,
	ChainOptimisticEthereum: chainOptimisticEthereumData,
	ChainBinanceSmartChain:  chainBscData,
	ChainOptimisticKovan:    chainOptimisticKovanData,
	ChainPolygon:            chainPolygonData,
}

var (
	chainEthereumData = &ChainData{
		id:        1,
		name:      "Ethereum",
		currency:  "ETH",
		jsonRPC:   "https://main-light.eth.linkpool.io/",
		wsRPC:     "wss://main-light.eth.linkpool.io/ws",
		blocktime: 13,
	}
	chainOptimisticEthereumData = &ChainData{
		id:       10,
		name:     "OptimisticEthereum",
		currency: "ETH",
		jsonRPC:  "https://mainnet.optimism.io/",
		wsRPC:    "https://ws-mainnet.optimism.io",
	}
	chainBscData = &ChainData{
		id:        56,
		name:      "Binance Smart Chain",
		currency:  "BNB",
		jsonRPC:   "https://bsc-dataseed1.defibit.io",
		wsRPC:     "wss://bsc-ws-node.nariox.org:443",
		blocktime: 3,
	}
	chainOptimisticKovanData = &ChainData{
		id:       69,
		name:     "OptimisticEthereum",
		currency: "ETH",
		jsonRPC:  "https://kovan.optimism.io",
		wsRPC:    "https://ws-kovan.optimism.io",
	}
	chainPolygonData = &ChainData{
		id:        137,
		name:      "Polygon",
		currency:  "MATIC",
		jsonRPC:   "https://rpc-mainnet.matic.network",
		wsRPC:     "wss://rpc-mainnet.matic.network",
		blocktime: 2,
	}
)
