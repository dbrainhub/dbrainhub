package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpClient interface {
	Send(ctx context.Context, url, method, body string) ([]byte, error)
}

func NewHttpClient(timeout time.Duration, retry int, retryInterval time.Duration) HttpClient {
	return &httpClientWithRetry{
		limit:    retry,
		interval: retryInterval,
		client: &httpClientWithTimeout{
			timeout: timeout,
			client:  &httpClient{},
		},
	}
}

type httpClientWithTimeout struct {
	timeout time.Duration
	client  HttpClient
}

func (h *httpClientWithTimeout) Send(ctx context.Context, url, method, body string) ([]byte, error) {
	ctxWithTiemout, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()

	return h.client.Send(ctxWithTiemout, url, method, body)
}

type httpClientWithRetry struct {
	limit    int
	interval time.Duration
	client   HttpClient
}

func (h *httpClientWithRetry) Send(ctx context.Context, url, method, body string) (resp []byte, err error) {
	for i := 0; i <= h.limit; i++ {
		resp, err = h.client.Send(ctx, url, method, body)
		if err == nil {
			return
		}
		time.Sleep(h.interval)
	}
	return
}

type httpClient struct {
	client http.Client
}

func (h *httpClient) Send(ctx context.Context, url, method, body string) ([]byte, error) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("new request error, err:%v", err)
	}

	request = request.WithContext(ctx)
	resp, err := h.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client do error, err:%v", err)
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp body error, err:%v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http resp code: %d, resp: %s", resp.StatusCode, string(respBytes))
	}
	return respBytes, nil
}
