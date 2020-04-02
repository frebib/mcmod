package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	Get = &cli.Command{
		Name:      "get",
		Usage:     "download a mod",
		Action:    cmdDoGet,
		ArgsUsage: "<name|id>",
	}
)

func cmdDoGet(c *cli.Context) (err error) {
	logger := c.App.Metadata["log"].(*logrus.Entry)
	log := logger.WithField("prefix", "get")

	if c.NArg() < 1 {
		log.Error("missing required arg: " + c.Command.ArgsUsage)
		return cli.ShowSubcommandHelp(c)
	}

	log.Infof("%s is '%s'", c.Command.ArgsUsage, c.Args().First())

	return nil
}

