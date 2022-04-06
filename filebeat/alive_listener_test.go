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
	listener := NewAliveListener("", httpClient, 0, 0)
	errCh := listener.Listen(context.Background())

	assert.Equal(t, (<-errCh).Error(), "111")
	assert.Equal(t, (<-errCh).Error(), "222")
	assert.Equal(t, (<-errCh).Error(), "333")

}
