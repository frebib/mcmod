package api

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	modlog "github.com/frebib/mcmod/log"
)

const ApiUrl = "https://addons-ecs.forgesvc.net/api"

const GameMinecraft int = 432

type ApiClient struct {
	HttpClient *http.Client
	ApiUrl     string
}

var DefaultClient = &ApiClient{
	HttpClient: new(http.Client),
	ApiUrl:     ApiUrl,
}

// ClientKey is the textual key used to identify
// an API Client inside a context.Context object
const ClientKey = "api-client"

// ClientFromContext loads a configured ApiClient from a context
// or returns the DefaultClient if none is found in the context
func ClientFromContext(ctx context.Context) *ApiClient {
	clientObj := ctx.Value(ClientKey)
	client, ok := clientObj.(*ApiClient)
	if ok {
		return client
	}
	return DefaultClient
}

func fetchJSON(ctx context.Context, client *http.Client, method, url string,
	body io.Reader) (*http.Response, error) {

	log := modlog.FromContext(ctx)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	log = log.WithField("request", fmt.Sprintf("%p", &req))

	log.Tracef("requesting %s %s", method, url)

	// All API responses return JSON
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	log.Tracef("request returned %s", resp.Status)

	// Non-200 responses are not useful, therefore error
	if resp.StatusCode != http.StatusOK {
		var bodyText string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err == nil {
			bodyText = string(bodyBytes)
		}
		return resp, &ErrHttpStatus{req, resp.StatusCode, bodyText}
	}
	return resp, err
}
func buildURL(base, urlPath string, params string) (string, error) {
	urlObj, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	path.Join()
	urlObj.Path = path.Join(urlObj.Path, urlPath)
	urlObj.RawQuery = params

	return urlObj.String(), nil
}
func buildURLParams(base, path string, params *url.Values) (string, error) {
	return buildURL(base, path, params.Encode())
}
