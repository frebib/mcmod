package download

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	modlog "github.com/frebib/mcmod/log"
	git_humanize "github.com/git-lfs/git-lfs/tools/humanize"
)

func FileFromURL(ctx context.Context, url, path string) error {
	rd, err := FromURL(ctx, nil, url)
	if err != nil {
		return err
	}
	return File(ctx, rd, path)
}

func File(ctx context.Context, src io.ReadCloser, path string) error {
	log := modlog.FromContext(ctx).
		WithField("name", path)

	srcCount, isCounter := src.(ReadCounter)
	if isCounter {
		log = log.WithField("size", humanize.IBytes(srcCount.ExpectedTotal()))
	}

	log.Info("downloading file")

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.WithError(err).Errorf("failed to create file")
		return err
	}

	start := time.Now()
	_, err = io.Copy(file, src)
	end := time.Now()
	if err != nil {
		log.WithError(err).Warnf("failed writing file")
		return err
	}
	// store the error and return it after the log entry
	err = src.Close()

	if isCounter {
		log = log.WithField("rate",
			git_humanize.FormatByteRate(srcCount.Count(), end.Sub(start)),
		)
	}
	log.Debugf("transferred in %s", end.Sub(start).Round(time.Millisecond))
	return err
}
