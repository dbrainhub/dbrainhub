package es

import (
	"fmt"
	"sync"

	"github.com/dbrainhub/dbrainhub/configs"
	esClient "github.com/elastic/go-elasticsearch/v8"
)

var client *esClient.Client
var once = sync.Once{}

func GetESClient() *esClient.Client {
	once.Do(func() {
		cfg := configs.GetGlobalServerConfig()
		esCfg := esClient.Config{
			Addresses: cfg.OutputServer.EsAddresses,
		}

		var err error
		client, err = esClient.NewClient(esCfg)
		if err != nil {
			panic(fmt.Sprintf("init es client error: %v", err))
		}
	})
	return client
}
