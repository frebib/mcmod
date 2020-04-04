package cmd

import (
	"context"
	"sort"

	"github.com/frebib/mcmod/api"
	modlog "github.com/frebib/mcmod/log"
)

type ModFilter struct {
	Release api.ReleaseType
	Version string
}

func listFilterMods(ctx context.Context, modID int, filter *ModFilter) (api.Files, error) {
	log := modlog.FromContext(ctx)

	files, err := api.ClientFromContext(ctx).Files(ctx, modID)
	if err != nil {
		log.WithError(err).Errorf("failed to list mod files")
		return nil, err
	}
	log.Debugf("found %d downloads", len(files))

	files, err = modFilter(ctx, files, filter)
	if err != nil {
		return nil, err
	}

	// Sort by release date, newest first
	sort.Sort(files)
	if len(files) < 1 {
		log.Warn("no download found")
		return nil, err
	}
	return files, nil
}

func modFilter(ctx context.Context, files api.Files, reqFilter *ModFilter) (api.Files, error) {
	log := modlog.FromContext(ctx)

	var filters []api.FileFilter
	if reqFilter.Release != api.ReleaseAny {
		log := log.WithField("release", reqFilter.Release.String())
		releaseFilter := api.FileFilterRelease(reqFilter.Release)
		releaseFilter.AfterFunc = func(files api.Files) error {
			log.Debugf("%d files match release filter", len(files))
			return nil
		}
		filters = append(filters, releaseFilter)
	}
	if reqFilter.Version != "" {
		log := log.WithField("version", reqFilter.Version)
		versionFilter := api.FileFilterVersion(reqFilter.Version)
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
