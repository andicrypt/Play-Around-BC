package listeners

import (
	"math/big"

	bridgeCore "github.com/axieinfinity/bridge-core"
	"github.com/ethereum/go-ethereum/log"
)

type GoerliTransaction interface {
	bridgeCore.Transaction
	Timestamp() uint64
}

func(g *GoerliListener) OrderCancelledHandler(fromChainId *big.Int, tx bridgeCore.Transaction, data []byte) error {
	log.Info("[OrderCancelledHandler] monitor")
	return nil
}

func(g *GoerliListener) OrderFulfilledHandler(fromChainId *big.Int, tx bridgeCore.Transaction, data []byte) error {
	log.Info("[OrderFulfilledHandler] monitor")
	return nil
}

func(g *GoerliListener) OrderValidatedHandler(fromChainId *big.Int, tx bridgeCore.Transaction, data []byte) error {
	log.Info("[OrderValidatedHandler] monitor")
	return nil
}