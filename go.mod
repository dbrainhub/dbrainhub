module github.com/dbrainhub/dbrainhub

go 1.16

replace github.com/Shopify/sarama => github.com/elastic/sarama v1.19.1-0.20220310193331-ebc2b0d8eef3

require (
	github.com/BurntSushi/toml v1.0.0
	github.com/gin-gonic/gin v1.7.2-0.20220325004437-205bb8151cb7
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/pkg/sftp v1.13.4
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/elastic/beats/v7 v7.17.2
	github.com/go-sql-driver/mysql v1.6.0
)
