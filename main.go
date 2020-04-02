package main

import (
	"context"
	"os"

	"github.com/frebib/mcmod/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	log     = logrus.WithContext(context.Background())
	lvlFlag = cli.StringFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Value:   logrus.InfoLevel.String(),
	}
)

func init() {
	log.Logger.SetOutput(os.Stderr)

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.QuoteCharacter = "'"
	log.Logger.Formatter = formatter
}

func main() {
	var app = &cli.App{
		Name:  "mcmod",
		Usage: "download, update and manage ad-hoc curseforge mod lists",

		HideHelpCommand:        true,
		UseShortOptionHandling: true,

		Commands: []*cli.Command{
			cmd.Get,
		},
		Flags: []cli.Flag{
			&lvlFlag,
		},
		Before: func(c *cli.Context) error {
			log := c.App.Metadata["log"].(*logrus.Entry)
			lvl, err := logrus.ParseLevel(c.String(lvlFlag.Name))
			if err != nil {
				return err
			}
			log.Logger.SetLevel(lvl)

			return nil
		},
		Metadata: map[string]interface{}{"log": log},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal()
	}
}
