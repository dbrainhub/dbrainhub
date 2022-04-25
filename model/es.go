package model

import (
	"bytes"
	"encoding/json"
	"fmt"

	esClient "github.com/elastic/go-elasticsearch/v8"
)

type EsClient interface {
	SendBatch(msgs []*ESMessage) error
}

func NewEsClient(addresses []string) (EsClient, error) {
	esCfg := esClient.Config{
		Addresses: addresses,
	}
	client, err := esClient.NewClient(esCfg)
	if err != nil {
		return nil, err
	}
	return &esClientImpl{
		client: client,
	}, nil
}

type esClientImpl struct {
	client *esClient.Client
}

func (e *esClientImpl) SendBatch(msgs []*ESMessage) error {
	if len(msgs) == 0 {
		return nil
	}
	var buf bytes.Buffer
	for _, msg := range msgs {
		meta, err := json.Marshal(msg.Meta)
		if err != nil {
			return fmt.Errorf("invalid meta, err: %v", err)
		}
		create := fmt.Sprintf(`{"create":%v}`, string(meta))

		data, err := json.Marshal(msg.Data)
		if err != nil {
			return fmt.Errorf("invalid data, err: %v", err)
		}

		buf.Grow(len(create) + len(data) + 2)
		buf.Write([]byte(create))
		buf.WriteByte('\n')
		buf.Write(data)
		buf.WriteByte('\n')
	}

	if _, err := e.client.Bulk(&buf); err != nil {
		return fmt.Errorf("es client bulk error, err: %v", err)
	}
	return nil
}

type ESMeta struct {
	Index    string `json:"_index"`
	Pipeline string `json:"pipeline,omitempty"`
}

type AgentIndexData struct {
	IP        string  `json:"ip"`
	Port      int     `json:"port"`
	CPURatio  float64 `json:"cpu_ratio"`
	MemRatio  float64 `json:"mem_ratio"`
	DiskRatio float64 `json:"disk_ratio"`
	QPS       float64 `json:"qps"`
	TPS       float64 `json:"tps"`
}

type ESMessage struct {
	Meta interface{}
	Data interface{}
}
