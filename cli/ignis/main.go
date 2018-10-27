package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tritonmedia/ignis/pkg/analysis"
	"github.com/tritonmedia/ignis/pkg/config"
	"github.com/tritonmedia/ignis/pkg/telegram"
	"github.com/tritonmedia/ignis/pkg/trello"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ignis"
	app.Usage = "start the ignis server"
	app.Author = "Jared Allard <jaredallard@outlook.com>"
	app.Action = func(c *cli.Context) error {
		d, err := os.Getwd()
		if err != nil {
			return err
		}

		err = analysis.Train()
		if err != nil {
			log.Printf("[main/analysis:train] failed to train model")
			os.Exit(1)
		}

		config, err := config.Load(filepath.Join(d, "config/config.yaml"))
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Config file not found in ./config/config.yaml ... :(")
				os.Exit(1)
			}
			return err
		}

		trello.Init(config.Trello.AppKey, config.Trello.Token, config.Trello.IDList)

		if config.Telegram.Token == "" {
			fmt.Println("Missing Telegram token in config.")
			os.Exit(1)
		}

		telegram.NewListener(config)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
