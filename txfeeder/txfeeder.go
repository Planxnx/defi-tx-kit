package txfeeder

import (
	"context"
	"log"
	"sync"

	"github.com/Planxnx/defi-tx-kit/contract"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
)

type LogsHandler func(types.Log) error

type TxFeeder struct {
	client               *ethclient.Client
	handlers             []func() (ethereum.Subscription, error)
	handlersSubscription []ethereum.Subscription
	handlersWg           *sync.WaitGroup
	closed               bool
	closedLock           sync.Mutex
}

func New(client *ethclient.Client) *TxFeeder {

	return &TxFeeder{
		client:     client,
		handlersWg: &sync.WaitGroup{},
	}
}

func (t *TxFeeder) Run() error {

	for _, handler := range t.handlers {
		sub, err := handler()
		if err != nil {
			return errors.Wrap(err, "Can't start handler")
		}
		t.handlersSubscription = append(t.handlersSubscription, sub)
	}

	go t.closeWhenAllHandlersStopped()

	return nil
}

func (t *TxFeeder) AddLogsListenr(ctx context.Context, eventFilter ethereum.FilterQuery, handler LogsHandler) error {

	t.handlers = append(t.handlers, func() (ethereum.Subscription, error) {
		txLogs := make(chan types.Log)
		sub, err := t.client.SubscribeFilterLogs(ctx, eventFilter, txLogs)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for txLog := range txLogs {
				t.handleTxLogs(txLog, handler)
			}
		}()
		return sub, err
	})

	return nil
}

func (t *TxFeeder) Close() error {
	t.closedLock.Lock()
	defer t.closedLock.Unlock()

	if t.closed {
		return nil
	}
	t.closed = true

	defer log.Println("TXFeeder client closed")

	for _, subscription := range t.handlersSubscription {
		subscription.Unsubscribe()
	}

	return nil
}

func (t *TxFeeder) handleTxLogs(txLogsData interface{}, handler interface{}) {

	//TODO: tracing or logging

	defer func() {
		if recovered := recover(); recovered != nil {
			log.Printf("Unexpected panic error: %+v\n", recovered)
		}
	}()

	//TODO: middleware

	var err error
	switch h := handler.(type) {
	case LogsHandler:
		err = h(txLogsData.(types.Log))
	case func(types.Log) error:
		err = h(txLogsData.(types.Log))
	case SwapHandlers:
		err = h(txLogsData.(*contract.PairSwap))
	case SyncHandlers:
		err = h(txLogsData.(*contract.PairSync))
	case TransferHandlers:
		err = h(txLogsData.(*contract.TokenTransfer))
	case ApprovalHandlers:
		err = h(txLogsData.(*contract.TokenApproval))
	default:
		panic(errors.Errorf("unknown txlogs handler type, data: %v", handler))
	}

	if err != nil {
		log.Printf("Error: handler error, %+v\n", err)
	}
}

func (t *TxFeeder) closeWhenAllHandlersStopped() {
	t.handlersWg.Wait()
	if t.isClosed() {
		return
	}

	log.Println("All txlogs handlers stopped, closing subcriber")

	t.Close()
}

func (t *TxFeeder) isClosed() bool {
	t.closedLock.Lock()
	defer t.closedLock.Unlock()

	return t.closed
}
