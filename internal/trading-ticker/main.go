package main

import (
	"context"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Planxnx/defi-tx-kit/contract"
	"github.com/Planxnx/defi-tx-kit/enums"
	"github.com/Planxnx/defi-tx-kit/txfeeder"
	"github.com/Planxnx/defi-tx-kit/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type LP struct {
	Address      common.Address
	Token0       *Token
	Token1       *Token
	Dex          string
	IsBaseToken0 bool
	Chain        utils.Chain
}

type Token struct {
	Symbol  string
	Address common.Address
	Chain   utils.Chain
}

var TokenData = map[string]*Token{
	"0xe9e7cea3dedca5984780bafc599bd69add087d56": {
		Symbol:  "BUSD",
		Address: common.HexToAddress("0xe9e7cea3dedca5984780bafc599bd69add087d56"),
		Chain:   utils.ChainBinanceSmartChain,
	},
	"0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c": {
		Symbol:  "BNB",
		Address: common.HexToAddress("0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"),
		Chain:   utils.ChainBinanceSmartChain,
	},
}

var lpData = map[string]*LP{
	"0x58f876857a02d6762e0101bb5c46a8c1ed44dc16": {
		Chain:        utils.ChainBinanceSmartChain,
		Address:      common.HexToAddress("0x58f876857a02d6762e0101bb5c46a8c1ed44dc16"),
		Token0:       TokenData["0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"],
		Token1:       TokenData["0xe9e7cea3dedca5984780bafc599bd69add087d56"],
		IsBaseToken0: true,
		Dex:          "Pancakeswap v2",
	},
	"0xaCAac9311b0096E04Dfe96b6D87dec867d3883Dc": {
		Chain:        utils.ChainBinanceSmartChain,
		Address:      common.HexToAddress("0xaCAac9311b0096E04Dfe96b6D87dec867d3883Dc"),
		Token0:       TokenData["0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"],
		Token1:       TokenData["0xe9e7cea3dedca5984780bafc599bd69add087d56"],
		IsBaseToken0: true,
		Dex:          "Biswap",
	},
}

func main() {

	ctx := context.Background()

	client, err := utils.ChainBinanceSmartChain.GetWSClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Event Filter
	queryPair := ethereum.FilterQuery{
		// Addresses: []common.Address{common.HexToAddress("0x58f876857a02d6762e0101bb5c46a8c1ed44dc16")},
	}

	txFeederClient := txfeeder.New(client)

	txFeederClient.AddSwapLogsHandler(ctx, queryPair, func(swapData *contract.PairSwap) error {

		lp, ok := lpData[strings.ToLower(swapData.Raw.Address.String())]
		if !ok {
			return nil
		}

		var side enums.SwapTxSide
		var swapRate *big.Float

		if lp.IsBaseToken0 {
			if swapData.Amount0Out.Sign() > 0 {
				side = enums.SwapBuy
			} else if swapData.Amount0In.Sign() > 0 {
				side = enums.SwapSell
			} else {
				log.Fatal("Error: cannot find swap side")
			}

			var amount0, amount1 *big.Float

			if side == enums.SwapBuy {
				amount0 = big.NewFloat(0).SetInt(swapData.Amount0Out)
				amount1 = big.NewFloat(0).SetInt(swapData.Amount1In)
			} else {
				amount0 = big.NewFloat(0).SetInt(swapData.Amount0In)
				amount1 = big.NewFloat(0).SetInt(swapData.Amount1Out)
			}

			swapRate = big.NewFloat(0).Quo(amount1, amount0)
		} else {
			if swapData.Amount1Out.Sign() > 0 {
				side = enums.SwapBuy
			} else if swapData.Amount1In.Sign() > 0 {
				side = enums.SwapSell
			} else {
				log.Fatal("Error: cannot find swap side")
			}

			var amount0, amount1 *big.Float

			if side == enums.SwapBuy {
				amount0 = big.NewFloat(0).SetInt(swapData.Amount0In)
				amount1 = big.NewFloat(0).SetInt(swapData.Amount1Out)
			} else {
				amount0 = big.NewFloat(0).SetInt(swapData.Amount0Out)
				amount1 = big.NewFloat(0).SetInt(swapData.Amount1In)
			}

			swapRate = big.NewFloat(0).Quo(amount0, amount1)
		}

		log.Printf("%v: %v %v\n", side, swapRate, lp.Token1.Symbol)

		return nil
	})

	log.Println("== Tx Logs Ticker ==")
	txFeederClient.Run()

	// Graceful shutdown
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
		syscall.SIGTERM, // kill -SIGTERM XXXX
	)

	<-signalChan
	log.Println("os.Interrupt - shutting down...")

	if err := txFeederClient.Close(); err != nil {
		log.Printf("Error while closing TXFeeder, %+v\n", err)
	}

	defer log.Println("gracefully stopped server")
}
