package main

import (
	"os"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "path to config file",
		},
	}
	app := &cli.App{
		Name:   "crawler",
		Usage:  "",
		Action: startCrawler,
		Flags:  flags,
		Commands: []cli.Command{
			{
				Name:   "listener",
				Usage:  "start listener which listens to contract events",
				Action: startListener,
				Flags:  flags,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}