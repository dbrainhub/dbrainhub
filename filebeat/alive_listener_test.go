package filebeat

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHttpClient struct {
	errStrs []string
	i       int
}

func (m *MockHttpClient) Send(ctx context.Context, url, method, body string) ([]byte, error) {
	err := fmt.Errorf(m.errStrs[m.i%len(m.errStrs)])
	m.i++
	return nil, err
}

func TestFilebeatAliveListener_Listen(t *testing.T) {
	httpClient := &MockHttpClient{
		errStrs: []string{"111", "222", "333"},
	}

	// use chan for testing logic in another gorounte
	errMsgs := make(chan string)
	callback := func(ctx context.Context, err error) {
		errMsgs <- err.Error()
	}
	listener := NewAliveListener("", httpClient, 0, 0, callback)
	listener.Listen(context.Background())

	assert.Equal(t, <-errMsgs, "111")
	assert.Equal(t, <-errMsgs, "222")
	assert.Equal(t, <-errMsgs, "333")

}
