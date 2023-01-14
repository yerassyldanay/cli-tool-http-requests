package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

// RequestMaker abstracts http requests
type RequestMaker interface {
	Do(*http.Request) (*http.Response, error)
}

// RequestHandler handles all urls provided by the user
type RequestHandler struct {
	// this allows to mock the methods of the struct without making real requests
	httpClient RequestMaker
	// when writing a test, you can pass zap.NewNop()
	logger *zap.Logger
}

func NewRequestHandler(logger *zap.Logger) *RequestHandler {
	// Create a custom HTTP client with timeouts
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return &RequestHandler{
		httpClient: httpClient,
		logger:     logger,
	}
}

// makeRequest makes http request to external service and sends the response back in the channel
func (rh RequestHandler) makeRequest(ctx context.Context, urlString string, responseCh chan<- Response) {
	u, err := url.Parse(urlString)
	if err != nil {
		responseCh <- Response{URL: urlString, Err: err}
		rh.logger.Error("invalid url", zap.String("url", urlString))
		return
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		responseCh <- Response{URL: urlString, Err: err}
		rh.logger.Error("error creating request", zap.String("url", urlString), zap.Error(err))
		return
	}

	// Add headers to the request
	req.Header.Add("User-Agent", "custom-agent")
	req.Header.Add("Accept-Encoding", "gzip")

	// Make the request
	resp, err := rh.httpClient.Do(req)
	if err != nil {
		responseCh <- Response{URL: urlString, Err: err}
		rh.logger.Debug("error making request", zap.String("url", urlString), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	// it is a good practice to read response body even if you dont use it
	// discard the response body by copying it to io.Discard
	contentLength, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		rh.logger.Error("failed to read response body", zap.String("url", urlString), zap.Error(err))
	}

	rh.logger.Debug("http response",
		zap.String("url", urlString),
		zap.String("status", resp.Status),
		zap.Int64("length", contentLength),
	)

	// Send the response through the channel
	responseCh <- Response{
		URL:     urlString,
		BodyLen: &contentLength,
		Err:     nil,
	}
}

// HandleUrls accepts list of urls and returns responses through the channel
func (rh RequestHandler) HandleUrls(ctx context.Context, urls []string) <-chan Response {
	// all the info will be sent to this channel
	responseCh := make(chan Response, 1)

	// this function is responsible for waiting all goroutines to have done their work
	// and closing the channel
	go func(responseCh chan Response, urls []string) {
		// close the channel when goroutines have finished sending responses
		// to the channel
		defer close(responseCh)

		// using wait group to sync all the goroutines
		var wg = &sync.WaitGroup{}
		for _, eachUrl := range urls {
			wg.Add(1)
			// this allows to add sync.WaitGroup without adding it to the makeRequest method
			go func(urlString string) {
				rh.makeRequest(ctx, urlString, responseCh)
				wg.Done()
			}(eachUrl)
		}
		wg.Wait()
	}(responseCh, urls)

	return responseCh
}
