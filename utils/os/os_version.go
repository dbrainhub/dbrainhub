package osutils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type VersionQuerier interface {
	// 发行版本
	GetOsVersion() (string, error)
	// 内核版本
	GetKernalVersion() (string, error)
}

func NewVersionQuerier(osType string) VersionQuerier {
	switch osType {
	case "darwin":
		return &darwinVersionQuerier{}
	case "linux":
		return &linuxVersionQuerier{}
	}
	return &unknownVersionQuerier{}
}

type darwinVersionQuerier struct {
}

func (d *darwinVersionQuerier) GetKernalVersion() (string, error) {
	const cmd = "uname -r"
	res, err := execCmd(cmd)
	return strings.TrimSpace(res), err
}

func (d *darwinVersionQuerier) GetOsVersion() (string, error) {
	const cmd = "sw_vers"
	const description = "ProductVersion:"
	output, err := execCmd(cmd)
	if err != nil {
		return "", err
	}
	// output eg:
	// ProductName:	Mac OS X
	// ProductVersion:	10.15.7
	// BuildVersion:	19H1519
	outputLines := strings.Split(output, "\n")
	if len(outputLines) < 2 {
		return "", fmt.Errorf("exec `%s` error, output: %s", cmd, output)
	}
	if !strings.Contains(outputLines[1], description) {
		return "", fmt.Errorf("unexpected output for exec `%s`, output: %s", cmd, output)
	}
	return strings.TrimSpace(strings.ReplaceAll(outputLines[1], description, "")), nil
}

type linuxVersionQuerier struct{}

func (d *linuxVersionQuerier) GetKernalVersion() (string, error) {
	const cmd = "uname -r"
	res, err := execCmd(cmd)
	return strings.TrimSpace(res), err
}

func (d *linuxVersionQuerier) GetOsVersion() (string, error) {
	const cmd = "lsb_release  -d"
	const description = "Description:"
	output, err := execCmd(cmd)
	if err != nil {
		return "", err
	}
	// output eg:
	// Description:	CentOS Linux release 7.6.1810 (Core)
	if !strings.Contains(output, description) {
		return "", fmt.Errorf("unexpected output for exec `%s`, output: %s", cmd, output)
	}
	return strings.TrimSpace(strings.ReplaceAll(output, description, "")), nil
}

type unknownVersionQuerier struct {
}

func (u *unknownVersionQuerier) GetKernalVersion() (string, error) {
	return "", nil
}

func (u *unknownVersionQuerier) GetOsVersion() (string, error) {
	return "", nil
}

func execCmd(cmd string) (string, error) {
	var buf bytes.Buffer
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = &buf
	err := c.Run()
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
