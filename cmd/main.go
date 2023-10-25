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
		Name:   "listener",
		Usage:  "start listener which listens to contract events",
		Flags:  flags,
		Action: startListener,
		Commands: []cli.Command{},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}