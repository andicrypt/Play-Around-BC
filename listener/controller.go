package listeners

import (
	"context"

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
	bridgeCore.AddListener("Goerli", InitEthereum)
	controller, err := bridgeCore.New(cfg, db, nil)

	if err != nil {
		panic(err)
	}

	return &Controller{controller}, nil
}

func InitEthereum(ctx context.Context, lsConfig *bridgeCore.LsConfig, store stores.MainStore, helpers utils.Utils, pool *bridgeCore.Pool) bridgeCore.Listener{
	goerliListener, err := NewGoerliListener(ctx, lsConfig, helpers, store, pool)
	if err != nil {
		log.Error("[GoerliListener]Error while init new ronin listener", "err", err)
		return nil
	}
	return goerliListener
}

func (c* Controller) Migrate(db *gorm.DB) error {
	return nil
}