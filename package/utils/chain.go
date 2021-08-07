package utils

import (
	"context"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	//Thanks Nethereum!
	infuraPrjectID = "7238211010344719ad14a89db874158c"
)

type Chain struct {
	// Network Settings
	id       int
	name     string
	currency string
	jsonRPC  string
	wsRPC    string

	// RPC Client
	ethclient *ethclient.Client
	once      sync.Once
}

var (
	ChainEthereum = Chain{
		id:       1,
		name:     "Ethereum",
		currency: "ETH",
		jsonRPC:  "https://mainnet.infura.io/v3/" + infuraPrjectID,
		wsRPC:    "wss://mainnet.infura.io/ws/v3/" + infuraPrjectID,
	}
	ChainOptimisticEthereum = Chain{
		id:       10,
		name:     "OptimisticEthereum",
		currency: "ETH",
		jsonRPC:  "https://mainnet.optimism.io/",
		wsRPC:    "https://ws-mainnet.optimism.io",
	}
	ChainBinanceSmartChain = Chain{
		id:       56,
		name:     "Binance Smart Chain",
		currency: "BNB",
		jsonRPC:  "https://bsc-dataseed1.defibit.io",
		wsRPC:    "wss://bsc-ws-node.nariox.org:443",
	}
	ChainOptimisticKovan = Chain{
		id:       69,
		name:     "OptimisticEthereum",
		currency: "ETH",
		jsonRPC:  "https://kovan.optimism.io",
		wsRPC:    "https://ws-kovan.optimism.io",
	}
	ChainPolygon = Chain{
		id:       137,
		name:     "Polygon",
		currency: "MATIC",
		jsonRPC:  "https://rpc-mainnet.matic.network",
		wsRPC:    "wss://rpc-mainnet.matic.network",
	}
)

func (c *Chain) ID() int {
	return c.id
}

func (c *Chain) Name() string {
	return c.name
}

func (c *Chain) Currency() string {
	return c.currency
}

func (c *Chain) GetClient(ctx context.Context) (client *ethclient.Client, err error) {
	if c.jsonRPC == "" {
		return nil, errors.New("undefined rpc endpoint")
	}
	c.once.Do(func() {
		c.ethclient, err = ethclient.DialContext(ctx, c.jsonRPC)
	})
	return c.ethclient, err
}

func (c *Chain) GetWSClient(ctx context.Context) (client *ethclient.Client, err error) {
	if c.wsRPC == "" {
		return nil, errors.New("undefined rpc endpoint")
	}
	c.once.Do(func() {
		c.ethclient, err = ethclient.DialContext(ctx, c.wsRPC)
	})
	return c.ethclient, err
}
