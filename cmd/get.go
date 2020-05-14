package cmd

import (
	"context"
	"sync"

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
			&flagNoDeps,
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

	var filelist []*api.File
	var filelock sync.Mutex

	var wait sync.WaitGroup
	var waitCh = make(chan struct{}, 1)
	var errs = make(chan error, 1)

	count := c.Args().Len()
	wait.Add(count)

	ctx, cancel := context.WithCancel(ctx)
	for _, mod := range c.Args().Slice() {
		go func(ctx context.Context, name string) {
			err := func(ctx context.Context, name string) error {
				defer wait.Done()

				mod, err := api.ClientFromContext(ctx).Lookup(ctx, name)
				if err != nil {
					return err
				}
				ctx, log = modlog.SetContextLogger(ctx, log.WithField("mod", mod.Slug))
				log.WithField("id", mod.ID).Info("found mod")

				filter := &ModFilter{Release: reqRelease, Version: reqVer}
				files, err := listFilterMods(ctx, mod.ID, filter)
				if err != nil {
					return err
				}

				// Pick the latest release
				modFile := files[0]
				log.WithField("file-id", modFile.ID).
					Tracef("chose '%s' as latest file", modFile.FileName)

				filelock.Lock()
				filelist = append(filelist, &modFile)
				filelock.Unlock()

				// Download optional dependencies, unless otherwise specified
				if !c.Bool(flagNoDeps.Name) {
					log.Debugf("resolving %d dependencies", len(modFile.Dependencies))
					var depwait sync.WaitGroup

					for _, dep := range modFile.Dependencies {
						depwait.Add(1)
						go func(ctx context.Context, depID int) {
							defer depwait.Done()
							// Update logger to display correct log details for the dep
							ctx, _ = modlog.SetContextLogger(ctx,
								log.WithFields(logrus.Fields{
									"mod":    depID,
									"dep-of": mod.Slug,
								}),
							)
							depMod, err := api.ClientFromContext(ctx).AddonByID(ctx, depID)
							if err != nil {
								log.WithError(err).Warnf("failed to lookup dependency")
							} else if depMod != nil {
								// Add the name now that we know what it is
								ctx, log = modlog.SetContextLogger(ctx,
									log.WithField("mod", depMod.Slug),
								)
							}

							depFiles, err := listFilterMods(ctx, depID, filter)
							if err != nil {
								return
							}
							if len(depFiles) > 0 {
								filelock.Lock()
								filelist = append(filelist, &depFiles[0])
								filelock.Unlock()
							} else {
								log.Warnf("no download found, skipping")
							}
						}(ctx, dep.AddonID)
					}
					depwait.Wait()
					log.Debugf("found an additional %d files", len(files)-1)
				}
				return nil
			}(ctx, name)
			if err != nil {
				errs <- err
			}
		}(ctx, mod)
	}

	log.Debug("wait for threads to end")

	go func() {
		wait.Wait()
		log.Debug("all threads finished")
		close(waitCh)
	}()

	select {
	case <-waitCh:
		break
	case err = <-errs:
		if err != nil {
			cancel()
			// Wait for all cancelled threads to exit
			wait.Wait()
			return err
		}
	}

	log.Debugf("identified %d mods, proceeding to download", count)

	for _, dl := range filelist {
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
