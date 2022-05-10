package ssh

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getSSHClient() (*SSHClient, error) {
	user := "ubuntu"
	addr := "192.168.64.3:22"

	priKey := `
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`

	pubKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCkMErajHMI9Xgf4T9HDbLQBZUGRkPRkRJt8RGWCJ2e6F6izNUjWSpAYv6wlENxe8xAxHJud8uA3KC8853CzKwkHIv/g4uTJbOcAYZYtqPD3IpSQ5nXYXRwjfpOBEbqnb/PhqC0XseN5I8fkxgQpMlqZ5LQKtf1IKGdaJ9ERxcpETO2DBI2qNWhA7lwih4y/eusU84MLCF9o10r1GWDar4iGymNZp1Njsc18X/xEuDswKoysWIswKr4wc8sFmmSQl2UlGysl0QOZLAborEiJA9gvfZKeOe0W1kDGFDV5tqXKLbxFzrTOEa9Bcc1b5DPqzr5YpoLtY/s2n4jKo54GKn9gHyZhsz/+2yli5Jk9jMr6U/B84otP+pQRfLGgFaCQaRFWhnWITBqqMLtmFNz9xTkqIqm/8bam1DMs+I8mNm2tipEMwl7qN1DBZs8iXMlvK5k1DMpAkyqPkjnc14AXEmZ3MbvT86r2HwB65m176NMowftA4r6nm4HBcU1lQiB3mM= ubuntu@primary`
	_ = pubKey

	return Connect(addr, user, []byte(priKey))
}

func TestSSH(t *testing.T) {
	t.Skip()
	c, err := getSSHClient()
	assert.Nil(t, err)

	defer c.Close()

	err = c.Scp("../../README.md", "README.md")
	assert.Nil(t, err)

	output, err := c.Run("cat README.md")
	assert.Nil(t, err)
	fmt.Printf("%s\n", string(output))

	_, err = c.Run("tar -czf README.tar.gz README.md")
	assert.Nil(t, err)

	err = c.Chmod("README.md", 0777)
	assert.Nil(t, err)
}
