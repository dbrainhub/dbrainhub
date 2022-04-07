package filebeat

import (
	"context"
	"fmt"
	"io"
)

// filebeat operation includes startup、reload.
// reload filebeat relys on model.FilebeatConfigModule conf. Ignore it in this interface.
type FilebeatOperation interface {
	Startup(ctx context.Context) (io.Reader, error)
}

type Executor interface {
	Exec(cmd string) (io.Reader, error)
}

func NewFilebeatOperation(exector Executor, executionFilepath, confFilePath, logPath, dataPath string) FilebeatOperation {
	return &filebeatOperationImpl{
		exector:           exector,
		executionFilePath: executionFilepath,
		confFilePath:      confFilePath,
		logPath:           logPath,
		dataPath:          dataPath,
	}
}

type filebeatOperationImpl struct {
	// some have defaul value.
	// details: https://www.elastic.co/guide/en/beats/filebeat/current/directory-layout.html
	exector           Executor
	executionFilePath string
	confFilePath      string
	logPath           string
	dataPath          string
}

func (f *filebeatOperationImpl) Startup(ctx context.Context) (io.Reader, error) {
	return f.exector.Exec(fmt.Sprintf("chmod +x %s; %s", f.executionFilePath, f.cmd()))
}

func (f *filebeatOperationImpl) cmd() string {
	res := fmt.Sprintf("%s -c %s", f.executionFilePath, f.confFilePath)
	if len(f.logPath) != 0 {
		res += fmt.Sprintf(" -path.logs %s", f.logPath)
	}
	if len(f.dataPath) != 0 {
		res += fmt.Sprintf(" -path.data %s", f.dataPath)
	}
	return res
}
