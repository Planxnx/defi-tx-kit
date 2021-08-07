package txfeeder

import (
	"context"

	"github.com/Planxnx/defi-tx-kit/package/contract"
	"github.com/Planxnx/defi-tx-kit/package/enums"
	"github.com/Planxnx/defi-tx-kit/package/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type TransferHandlers func(*contract.TokenTransfer) error

type ApprovalHandlers func(*contract.TokenApproval) error

func (t *TxFeeder) AddTransferLogsHandler(ctx context.Context, eventFilter ethereum.FilterQuery, handler TransferHandlers) error {

	eventFilter.Topics = [][]common.Hash{
		{utils.TokenABI.Events[enums.TokenEventTransfer.ToString()].ID},
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
					data, err := utils.TokenParseTransfer(txLog)
					if err != nil {
						return errors.Wrap(err, "Can't parse transfer logs data")
					}
					return handler(data)
				})
			}
		}()

		return sub, nil
	})

	return nil
}

func (t *TxFeeder) AddApprovalLogsHandler(ctx context.Context, eventFilter ethereum.FilterQuery, handler ApprovalHandlers) error {

	eventFilter.Topics = [][]common.Hash{
		{utils.TokenABI.Events[enums.TokenEventApproval.ToString()].ID},
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
					data, err := utils.TokenParseApproval(txLog)
					if err != nil {
						return errors.Wrap(err, "Can't parse approval logs data")
					}
					return handler(data)
				})
			}
		}()

		return sub, nil
	})

	return nil
}
