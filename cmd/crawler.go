package main

import (

	"github.com/andicrypt/play-around-bc/ccip_reader"
	"github.com/andicrypt/play-around-bc/common"
	"github.com/andicrypt/play-around-bc/handler"
	"gopkg.in/urfave/cli.v1"
)

type config struct {
	Database  *common.Database         `json:"database" mapstructure:"database"`
	Server    *ccip_reader.Config      `json:"server" mapstructure:"server"`
	CDHandler *handler.Config `json:"ens" mapstructure:"ens"`
	Testing   bool                     `json:"testing" mapstructure:"testing"`
	Domain    *Domain                  `json:"domain" mapstructure:"domain"`
}

func startCrawler(ctx *cli.Context) {
	cfg := &config{}
	common.Load(ctx.String("config"), cfg)

	db, err := common.NewDBConn(cfg.Database, cfg.Testing)
	if err != nil {
		panic(err)
	}

	// new ccip server
	server := ccip_reader.NewServer(cfg.Server)
	handler, err := handler.NewDCHandler(cfg.CDHandler, db)
	if err != nil {
		panic(err)
	}
	server.Add(handler.GetHandlers())
	if err = server.Start(); err != nil {
		panic (err)
	}
}