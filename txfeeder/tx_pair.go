package txfeeder

import (
	"context"

	"github.com/Planxnx/defi-tx-kit/contract"
	"github.com/Planxnx/defi-tx-kit/enums"
	"github.com/Planxnx/defi-tx-kit/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type SwapHandlers func(*contract.PairSwap) error

type SyncHandlers func(*contract.PairSync) error

type SwapAndSyncHandlers struct {
	SwapHandler SwapHandlers
	SyncHandler SyncHandlers
}

func (t *TxFeeder) AddSwapLogsHandler(ctx context.Context, eventFilter ethereum.FilterQuery, handler SwapHandlers) error {

	eventFilter.Topics = [][]common.Hash{
		{utils.PairABI.Events[enums.PairEventSwap.ToString()].ID},
	}

	t.handlers = append(t.handlers, func() (ethereum.Subscription, error) {
		txLogs := make(chan types.Log)
		sub, err := t.client.SubscribeFilterLogs(ctx, eventFilter, txLogs)
		if err != nil {
			return nil, errors.Wrap(err, "Start subscribe error")
		}

		t.handlersWg.Add(1)
		go func() {
			defer t.handlersWg.Done()

			for txLog := range txLogs {
				t.handleTxLogs(txLog, func(txLog types.Log) error {
					data, err := utils.PairParseSwap(txLog)
					if err != nil {
						return errors.Wrap(err, "Can't parse swap logs data")
					}
					return handler(data)
				})
			}
		}()

		return sub, nil
	})

	return nil
}

func (t *TxFeeder) AddSyncLogsHandler(ctx context.Context, eventFilter ethereum.FilterQuery, handler SyncHandlers) error {

	eventFilter.Topics = [][]common.Hash{
		{utils.PairABI.Events[enums.PairEventSync.ToString()].ID},
	}

	t.handlers = append(t.handlers, func() (ethereum.Subscription, error) {
		txLogs := make(chan types.Log)
		sub, err := t.client.SubscribeFilterLogs(ctx, eventFilter, txLogs)
		if err != nil {
			return nil, errors.Wrap(err, "Start subscribe error")
		}

		t.handlersWg.Add(1)
		go func() {
			defer t.handlersWg.Done()

			for txLog := range txLogs {
				t.handleTxLogs(txLog, func(txLog types.Log) error {
					data, err := utils.PairParseSync(txLog)
					if err != nil {
						return errors.Wrap(err, "Can't parse sync logs data")
					}
					return handler(data)
				})
			}
		}()

		return sub, nil
	})

	return nil
}

func (t *TxFeeder) AddSwapAndSyncLogsHandler(ctx context.Context, eventFilter ethereum.FilterQuery, handlers *SwapAndSyncHandlers) error {

	eventFilter.Topics = [][]common.Hash{
		{utils.PairABI.Events[enums.PairEventSync.ToString()].ID, utils.PairABI.Events[enums.PairEventSwap.ToString()].ID},
	}

	t.handlers = append(t.handlers, func() (ethereum.Subscription, error) {
		txLogs := make(chan types.Log)
		sub, err := t.client.SubscribeFilterLogs(ctx, eventFilter, txLogs)
		if err != nil {
			return nil, errors.Wrap(err, "Start subscribe error")
		}

		t.handlersWg.Add(1)
		go func(messages <-chan types.Log) {
			defer t.handlersWg.Done()

			for txLog := range txLogs {
				topic := txLog.Topics[0]
				switch {
				case topic == utils.PairABI.Events[enums.PairEventSync.ToString()].ID && handlers.SyncHandler != nil:
					t.handleTxLogs(txLog, func(txLog types.Log) error {
						data, err := utils.PairParseSync(txLog)
						if err != nil {
							return errors.Wrap(err, "Can't parse sync logs data")
						}
						return handlers.SyncHandler(data)
					})

				case topic == utils.PairABI.Events[enums.PairEventSwap.ToString()].ID && handlers.SwapHandler != nil:
					t.handleTxLogs(txLog, func(txLog types.Log) error {
						data, err := utils.PairParseSwap(txLog)
						if err != nil {
							return errors.Wrap(err, "Can't parse swap logs data")
						}
						return handlers.SwapHandler(data)
					})
				}
			}
		}(txLogs)

		return sub, nil
	})

	return nil
}
