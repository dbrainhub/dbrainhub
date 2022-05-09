package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	esClient "github.com/elastic/go-elasticsearch/v8"
)

type (
	IndexQuerier interface {
		Query(ctx context.Context, param *IndexQuerierParam) ([]*IndexQuerierResult, error)
	}

	IndexQuerierParam struct {
		Index    string
		Begin    time.Time
		End      time.Time
		Buckets  int64
		AggField string
		Cond     map[string]interface{}
	}

	IndexQuerierResult struct {
		Count     int
		Value     float64
		StartTime string
	}
)

func NewIndexQuerier(client *esClient.Client) IndexQuerier {
	return &indexQuerier{
		client: client,
	}
}

type indexQuerier struct {
	client *esClient.Client
}

func (i *indexQuerier) Query(ctx context.Context, param *IndexQuerierParam) ([]*IndexQuerierResult, error) {
	from := param.Begin.Format("2006-01-02T15:04:05.999Z")
	to := param.End.Format("2006-01-02T15:04:05.999Z")
	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"avgs": map[string]interface{}{ // `avgs` is a custom string
				"date_histogram": genHistogram(from, to, param.Buckets),
				"aggs": map[string]interface{}{
					"avg_value": map[string]interface{}{ // `avg_value` is a custom string
						"avg": map[string]interface{}{
							"field": param.AggField,
						},
					},
				},
			},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []map[string]interface{}{
					{
						"bool": map[string]interface{}{
							"filter": genTerms(param.Cond),
						},
					},
					{
						"range": genRange(from, to),
					},
				},
			},
		},
		"size": 0,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("format es query body err: %v", err)
	}

	searchRes, err := i.client.Search(
		i.client.Search.WithContext(ctx),
		i.client.Search.WithIndex(param.Index),
		i.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("query es error, err: %v, body: %s", err, buf.String())
	}
	defer searchRes.Body.Close()

	var res indexQuerierResponse
	if err := json.NewDecoder(searchRes.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("parse es search result err: %v, body: %s", err, buf.String())
	}

	return res.ToResult(), nil
}

type indexQuerierResponse struct {
	Aggregations struct {
		Avgs struct {
			Buckets []*indexQuerierBucket `json:"buckets"`
		} `json:"avgs"`
	} `json:"aggregations"`
}

type indexQuerierBucket struct {
	KeyAsString string `json:"key_as_string"`
	AvgValue    struct {
		Value float64 `json:"value"`
	} `json:"avg_value"`
	DocCount int `json:"doc_count"`
}

func (a *indexQuerierResponse) ToResult() []*IndexQuerierResult {
	var res []*IndexQuerierResult
	for _, bucket := range a.Aggregations.Avgs.Buckets {
		res = append(res, &IndexQuerierResult{
			Count:     bucket.DocCount,
			Value:     bucket.AvgValue.Value,
			StartTime: bucket.KeyAsString,
		})
	}
	return res
}
