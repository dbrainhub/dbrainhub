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

func (m *MockExecutor) Exec(cmd string) (io.Reader, error) {
	return strings.NewReader(cmd), nil
}

func TestFilebeatStartup(t *testing.T) {
	filebeatOp := NewFilebeatOperation(&MockExecutor{},
		"/usr/share/filebeat/bin/filebeat",
		"/etc/filebeat.yml",
		"",
		"/var/log/filebeat")

	res, err := filebeatOp.Startup(context.Background())
	assert.Nil(t, err)

	cmd := make([]byte, 1024)
	n, err := res.Read(cmd)
	assert.Nil(t, err)
	assert.NotEqual(t, n, 0)
	assert.Equal(t, string(cmd[:n]), "chmod +x /usr/share/filebeat/bin/filebeat; /usr/share/filebeat/bin/filebeat -c /etc/filebeat.yml -path.data /var/log/filebeat")
}
