package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	bridge_core "github.com/axieinfinity/bridge-core"
	"github.com/axieinfinity/bridge-core/models"
	"github.com/axieinfinity/bridge-core/stores"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"

	"github.com/andicrypt/play-around-bc/common"
	listeners "github.com/andicrypt/play-around-bc/listener"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultMaxLogsBatch   = 100
	defaultSafeBlockRange = 41926
	defaultRPC = "https://eth-goerli.g.alchemy.com/v2/0EcWWGpuYHAOuXQhHeagpC_UhlEXvfYh"
	defaultFromBlock = 9023330

)


type Listener struct {
	RPC            string               `yaml:"rpc" mapstructure:"rpc"`
	FromBlock      uint64               `yaml:"from_block" mapstructure:"from_block"`
	Handlers       []*Handler           `yaml:"handlers" mapstructure:"handlers"`
	Contracts      map[string]*Contract `yaml:"contracts" mapstructure:"contracts"`
	SafeBlockRange uint64               `yaml:"safe_block_range" mapstructure:"safe_block_range"`
	LogBatchSize   int                  `yaml:"log_batch_size" mapstructure:"log_batch_size"`
}

type Domain struct {
	Ronin    string `yaml:"ronin" mapstructure:"ronin"`
	Ethereum string `yaml:"ethereum" mapstructure:"ethereum"`
}

type Contract struct {
	Address string `yaml:"address" mapstructure:"address"`
	AbiPath string `yaml:"abi" mapstructure:"abi"`
	Name    string `yaml:"name" mapstructure:"name"`
}

type Handler struct {
	Contract    string `yaml:"contract" mapstructure:"contract"`
	Event       string `yaml:"event" mapstructure:"event"`
	Handler     string `yaml:"handler" mapstructure:"handler"`
	Description string `yaml:"description" mapstructure:"description"`
}

type Workers struct {
	NumberOfWorkers int   `yaml:"numberOfWorkers" mapstructure:"numberOfWorkers"`
	MaxQueueSize    int   `yaml:"maxQueueSize" mapstructure:"maxQueueSize"`
	MaxRetry        int32 `yaml:"maxRetry" mapstructure:"maxRetry"`
	BackOff         int32 `yaml:"backoff" mapstructure:"backoff"`
}

type ListenerConfig struct {
	ChainId   string            `yaml:"chainId" mapstructure:"chainId"`
	Listener  *Listener         `yaml:"listener" mapstructure:"listener"`
	HTTP      *common.HTTP      `yaml:"http" mapstructure:"http"`
	Workers   *Workers          `yaml:"workers" mapstructure:"workers"`
	Verbosity int               `yaml:"verbosity" mapstructure:"verbosity"`
	DB        *common.Database  `yaml:"database" mapstructure:"database"`
	Testing   bool              `yaml:"testing" mapstructure:"testing"`
	Secrets   map[string]string `yaml:"secrets" mapstructure:"secrets"`
	Domain    *Domain           `yaml:"domain" mapstructure:"domain"`
}

func loadConfigAndDB(path string) (*bridge_core.Config, *gorm.DB) {
	fmt.Println("kakaka", path)
	cfg := &ListenerConfig{}
	common.Load(path, cfg)
	fmt.Printf("my config2 %+v\n", cfg)
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(cfg.Verbosity), log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	if cfg.Listener.SafeBlockRange == 0 {
		cfg.Listener.SafeBlockRange = defaultSafeBlockRange
	}
	if cfg.Listener.RPC == "" {
		cfg.Listener.RPC = defaultRPC
	}
	if cfg.Listener.SafeBlockRange == 0 {
		cfg.Listener.SafeBlockRange = defaultSafeBlockRange
	}
	if cfg.Listener.FromBlock == 0 {
		cfg.Listener.FromBlock = defaultFromBlock
	}

	// load ethereum listener
	ethereumConfig := &bridge_core.LsConfig{
		ChainId:          cfg.ChainId,
		Name:             "ETHEREUM",
		RpcUrl:           cfg.Listener.RPC,
		SafeBlockRange:   cfg.Listener.SafeBlockRange,
		FromHeight:       cfg.Listener.FromBlock,
		Subscriptions:    make(map[string]*bridge_core.Subscribe),
		GetLogsBatchSize: defaultMaxLogsBatch,
		Contracts:        make(map[string]string),
	}

	// ethereumConfig := &bridge_core.LsConfig{
	// 	ChainId:          cfg.ChainId,
	// 	Name:             "ETHEREUM",
	// 	RpcUrl:           "https://eth-goerli.g.alchemy.com/v2/0EcWWGpuYHAOuXQhHeagpC_UhlEXvfYh",
	// 	SafeBlockRange:   3,
	// 	FromHeight:       7126945,
	// 	Subscriptions:    make(map[string]*bridge_core.Subscribe),
	// 	GetLogsBatchSize: defaultMaxLogsBatch,
	// 	Contracts:        make(map[string]string),
	// }

	if cfg.Listener.LogBatchSize > 0 {
		ethereumConfig.GetLogsBatchSize = cfg.Listener.LogBatchSize
	}
	for name, contract := range cfg.Listener.Contracts {
		ethereumConfig.Contracts[name] = contract.Address
	}

	for _, handler := range cfg.Listener.Handlers {
		// load contract info
		contract, ok := cfg.Listener.Contracts[handler.Contract]
		if !ok {
			panic(fmt.Sprintf("%s not found", handler.Contract))
		}
		// load abi from ABI path
		fmt.Println("ppppp", contract.AbiPath)
		abiFile, err := os.ReadFile(contract.AbiPath)
		if err != nil {
			panic(err)
		}
		smcAbi, err := abi.JSON(bytes.NewBuffer(abiFile))
		if err != nil {
			panic(err)
		}
		ethereumConfig.Subscriptions[handler.Event] = &bridge_core.Subscribe{
			To:   contract.Address,
			Type: 1,
			Handler: &bridge_core.Handler{
				Name: handler.Event,
				ABI:  &smcAbi,
			},
			CallBacks: map[string]string{
				"ethereum": handler.Handler,
			},
		}
	}

	bridgeConfig := &bridge_core.Config{
		Listeners: map[string]*bridge_core.LsConfig{
			"ethereum": ethereumConfig,
		},
		DB: &stores.Database{
			Host:            cfg.DB.Host,
			User:            cfg.DB.User,
			Password:        cfg.DB.Password,
			DBName:          cfg.DB.DBName,
			Port:            int(cfg.DB.Port),
			ConnMaxLifetime: cfg.DB.ConnMaxLifetime,
			MaxIdleConns:    cfg.DB.MaxIdleConns,
			MaxOpenConns:    cfg.DB.MaxOpenConns,
		},
		Testing: cfg.Testing,
	}

	if cfg.Workers != nil {
		if cfg.Workers.NumberOfWorkers > 0 {
			bridgeConfig.NumberOfWorkers = cfg.Workers.NumberOfWorkers
		}
		if cfg.Workers.MaxQueueSize > 0 {
			bridgeConfig.MaxQueueSize = cfg.Workers.MaxQueueSize
		}
		if cfg.Workers.BackOff > 0 {
			bridgeConfig.BackOff = cfg.Workers.BackOff
		}
		if cfg.Workers.MaxRetry > 0 {
			bridgeConfig.MaxRetry = cfg.Workers.MaxRetry
		}
	}

	db, err := common.NewDBConn(cfg.DB, cfg.Testing)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.ProcessedBlock{},
		&models.Job{},
	)

	return bridgeConfig, db
}

func startListener(ctx *cli.Context) {
	fmt.Println("hello 0")
	fmt.Println("hhhh", ctx.String("config"))

	cfg, db := loadConfigAndDB(ctx.String("config"))
	// cfg, db := loadConfigAndDB(path)

	fmt.Printf("config10 %+v", cfg)
	controller, err := listeners.NewController(cfg, db, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("hello 4")

	fmt.Printf("controller %+v", controller)
	if err = controller.Start(); err != nil {
		panic(err)
	}
	fmt.Println("hello 5")
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)
	<-sigc
	go func() {
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				log.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
		debug.SetTraceback("all")
		panic("boom")
	}()
	controller.Close()
}