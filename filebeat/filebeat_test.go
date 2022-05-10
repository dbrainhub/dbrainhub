package filebeat

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockExecutor struct {
}

func (m *MockExecutor) Exec(ctx context.Context, cmd string) (io.Reader, io.Reader, error) {
	return strings.NewReader(cmd), nil, nil
}

func TestFilebeatStartup(t *testing.T) {
	filebeatOp := NewFilebeatOperation(&MockExecutor{},
		"/usr/share/filebeat/bin/filebeat",
		"/etc/filebeat.yml",
		"/usr/share/filebeat/bin/",
		"",
		"/var/log/filebeat")

	res, _, err := filebeatOp.Startup(context.Background())
	assert.Nil(t, err)

	cmd := make([]byte, 1024)
	n, err := res.Read(cmd)
	assert.Nil(t, err)
	assert.NotEqual(t, n, 0)
	assert.Equal(t, string(cmd[:n]), "chmod +x /usr/share/filebeat/bin/filebeat && cd /usr/share/filebeat/bin/ && /usr/share/filebeat/bin/filebeat -c /etc/filebeat.yml -path.data /var/log/filebeat")
}
