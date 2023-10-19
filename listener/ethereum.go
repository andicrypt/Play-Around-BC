package listeners

import (
	"context"
	"fmt"

	bridgeCore "github.com/axieinfinity/bridge-core"
	"github.com/axieinfinity/bridge-core/stores"
	bridgev2_listener "github.com/axieinfinity/bridge-v2/listener"
	bridgeCoreUtils "github.com/axieinfinity/bridge-core/utils"
	"github.com/ethereum/go-ethereum/log"
)



type GoerliListener struct {
	*bridgev2_listener.EthereumListener
}


func NewGoerliListener(ctx context.Context, cfg *bridgeCore.LsConfig, helpers bridgeCoreUtils.Utils, store stores.MainStore, pool *bridgeCore.Pool) (*GoerliListener, error) {
	ethereumListener, err := bridgev2_listener.NewEthereumListener(ctx, cfg, helpers, store, pool)
	if err != nil {
		log.Error(fmt.Sprintf("[New%sListener] error while initialize Ethereum Listener", cfg.Name), "err", err, "url", cfg.RpcUrl)
		return nil, err
	}

	goerliListener := &GoerliListener{ethereumListener}

	return goerliListener, nil
}



