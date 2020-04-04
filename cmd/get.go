package cmd

import (
	"sort"

	"github.com/frebib/mcmod/api"
	"github.com/frebib/mcmod/download"
	modlog "github.com/frebib/mcmod/log"
	"github.com/frebib/mcmod/util"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	Get = &cli.Command{
		Name:      "get",
		Usage:     "download a mod",
		Action:    cmdDoGet,
		ArgsUsage: "<name|id>",
		Flags: []cli.Flag{
			&flagDirectory,
			&flagOutputFile,
			&flagRelease,
			&flagVersion,
		},
	}
)

func cmdDoGet(c *cli.Context) (err error) {
	ctx := c.Context
	log := modlog.FromContext(ctx)

	if c.NArg() < 1 {
		log.Error("missing required arg: " + c.Command.ArgsUsage)
		return cli.ShowSubcommandHelp(c)
	}

	reqVer := c.String("gamever")
	reqReleaseText := c.String("release")
	reqRelease := api.ParseReleaseType(reqReleaseText)
	if reqRelease == api.ReleaseUnknown {
		log.Errorf("invalid release type '%s'", reqReleaseText)
		return
	}

	mod, err := api.DefaultClient.Lookup(ctx, c.Args().First())
	if err != nil {
		return err
	}
	log = log.WithFields(logrus.Fields{"id": mod.ID, "name": mod.Slug})
	log.Info("found mod")

	files, err := api.DefaultClient.Files(ctx, mod.ID)
	if err != nil {
		log.WithError(err).Errorf("failed to list mod files")
		return err
	}
	log.Debugf("found %d downloads", len(files))

	files, err = modFilter(ctx, files, reqRelease, reqVer)
	if err != nil {
		return err
	}

	// Sort by release date, newest first
	sort.Sort(files)
	if len(files) < 1 {
		log.Warn("no download found")
		return nil
	}

	// Pick the latest release
	var toDownload = []api.File{files[0]}

	for _, dl := range toDownload {
		// Calculate final path+filename for mod output
		outFile := c.String(flagOutputFile.Name)
		outDir := c.String(flagDirectory.Name)
		filePath, err := util.CalcFilePath(dl.FileName, outFile, outDir)
		if err != nil {
			return err
		}

		err = download.FileFromURL(ctx, dl.DownloadURL, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
