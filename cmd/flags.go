package cmd

import (
	"github.com/frebib/mcmod/api"
	"github.com/urfave/cli/v2"
)

var (
	flagDirectory = cli.PathFlag{
		Name:        "directory",
		Usage:       "output directory for downloaded file",
		DefaultText: "$PWD",
		Aliases:     []string{"d"},
		EnvVars:     []string{"OUTPUT_DIRECTORY"},
	}
	flagOutputFile = cli.PathFlag{
		Name:    "output",
		Usage:   "output filename",
		Aliases: []string{"o"},
		EnvVars: []string{"OUTPUT_FILENAME"},
	}
	flagVersion = cli.StringFlag{
		Name:    "gamever",
		Usage:   "game version",
		Aliases: []string{"V"},
		EnvVars: []string{"MINECRAFT_VERSION"},
	}
	flagRelease = cli.StringFlag{
		Name:    "release",
		Usage:   "release type, of [any, release, beta, alpha]",
		Aliases: []string{"r"},
		Value:   api.ReleaseRelease.String(),
		EnvVars: []string{"RELEASE_TYPE"},
	}
)
