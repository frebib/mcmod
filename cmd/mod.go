package cmd

import (
	"context"

	"github.com/frebib/mcmod/api"
	modlog "github.com/frebib/mcmod/log"
)

func modFilter(ctx context.Context, files api.Files, reqRelease api.ReleaseType, reqVer string) (api.Files, error) {
	log := modlog.FromContext(ctx)

	var filters []api.FileFilter
	if reqRelease != api.ReleaseAny {
		log := log.WithField("release", reqRelease.String())
		releaseFilter := api.FileFilterRelease(reqRelease)
		releaseFilter.AfterFunc = func(files api.Files) error {
			log.Debugf("%d files match release filter", len(files))
			return nil
		}
		filters = append(filters, releaseFilter)
	}
	if reqVer != "" {
		log := log.WithField("version", reqVer)
		versionFilter := api.FileFilterVersion(reqVer)
		versionFilter.AfterFunc = func(files api.Files) error {
			log.Debugf("%d files match version filter", len(files))
			return nil
		}
		filters = append(filters, versionFilter)
	}

	// Apply requested filters
	if len(filters) > 0 {
		var err error
		files, err = files.Filter(filters)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
}
