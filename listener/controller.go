package listeners

import (
	"context"
	"fmt"

	bridgeCore "github.com/axieinfinity/bridge-core"
	"github.com/axieinfinity/bridge-core/stores"
	"github.com/axieinfinity/bridge-core/utils"
	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)
type Controller struct {
	*bridgeCore.Controller
}

func NewController(cfg *bridgeCore.Config, db *gorm.DB, helpers utils.Utils) (*Controller, error) {
	bridgeCore.AddListener("ethereum", InitEthereum)
	fmt.Println("Debug 2222")
	controller, err := bridgeCore.New(cfg, db, helpers)

	if err != nil {
		panic(err)
	}

	return &Controller{controller}, nil
}

func InitEthereum(ctx context.Context, lsConfig *bridgeCore.LsConfig, store stores.MainStore, helpers utils.Utils, pool *bridgeCore.Pool) bridgeCore.Listener{
	fmt.Println("Debug abcd")

	goerliListener, err := NewGoerliListener(ctx, lsConfig, helpers, store, pool)
	if err != nil {
		log.Error("[GoerliListener]Error while init new ethereum listener", "err", err)
		return nil
	}
	log.Info("Finished initializing Ethereum listener")
	fmt.Println("Debug 5")
	fmt.Printf("goerliListener: %+v", goerliListener)

	return goerliListener
}
