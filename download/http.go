package download

import (
	"context"
	"net/http"
	"strconv"

	"github.com/frebib/mcmod/api"
	modlog "github.com/frebib/mcmod/log"
)

func FromURL(ctx context.Context, client *http.Client, url string) (ReadCounter, error) {
	log := modlog.FromContext(ctx)

	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = &api.ErrHttpStatus{Req: resp.Request, Code: resp.StatusCode}
		log.WithError(err).
			WithField("status", resp.StatusCode).
			Errorf("got download error")
	}

	totalStr := resp.Header.Get("content-length")
	total, err := strconv.ParseUint(totalStr, 10, 64)
	if err != nil {
		total = 0
	}

	return &CountingReader{Reader: resp.Body, Total: total}, err
}
