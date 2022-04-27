package model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dbrainhub/dbrainhub/utils/search_time"
)

type (
	EsAvgAggs interface {
		AggregateAvg(ctx context.Context, field string, beginTime, endTime time.Time, cond map[string]interface{}) ([]*AvgAggsResult, error)
	}

	AvgAggsResult struct {
		Count     int
		Value     float64
		StartTime string
	}
)

func (e *esClientImpl) AggregateAvg(ctx context.Context, field string, beginTime, endTime time.Time, cond map[string]interface{}) ([]*AvgAggsResult, error) {
	from := beginTime.Format("2006-01-02T15:04:05.999Z")
	to := endTime.Format("2006-01-02T15:04:05.999Z")
	query := map[string]interface{}{
		"aggs": map[string]interface{}{
			"avgs": map[string]interface{}{
				"date_histogram": map[string]interface{}{
					"field":          "@timestamp",
					"fixed_interval": search_time.GetInterval(from, to, 50),
					"min_doc_count":  1, // ignore count == 0
					"format":         "yyyy-MM-dd HH:mm:ss",
				},
				"aggs": map[string]interface{}{
					"avg_value": map[string]interface{}{
						"avg": map[string]interface{}{
							"field": field,
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
							"filter": genTerms(cond),
						},
					},
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"format": "strict_date_optional_time",
								"gte":    from,
								"lte":    to,
							},
						},
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

	fmt.Println(buf.String())
	searchRes, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("test-index"),
		e.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("query es error, err: %v, body: %s", err, buf.String())
	}
	defer searchRes.Body.Close()

	var res avgAggsResponse
	if err := json.NewDecoder(searchRes.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("parse es search result err: %v, body: %s", err, buf.String())
	}

	return res.ToAvgAggsResults(), nil
}

func genTerms(cond map[string]interface{}) []map[string]interface{} {
	res := []map[string]interface{}{}
	for key, value := range cond {
		res = append(res, map[string]interface{}{
			"term": map[string]interface{}{
				key: value,
			},
		})
	}
	/*
		res = append(res, map[string]interface{}{
			"range": map[string]interface{}{
				"@timestamp": map[string]interface{}{
					"format": "strict_date_optional_time",
					"gte":    from,
					"lte":    to,
				},
			}})
	*/
	return res
}

type avgAggsResponse struct {
	Aggregations struct {
		Avgs struct {
			Buckets []*avgAggsBucket `json:"buckets"`
		} `json:"avgs"`
	} `json:"aggregations"`
}

type avgAggsBucket struct {
	KeyAsString string `json:"key_as_string"`
	AvgValue    struct {
		Value float64 `json:"value"`
	} `json:"avg_value"`
	DocCount int `json:"doc_count"`
}

func (a *avgAggsResponse) ToAvgAggsResults() []*AvgAggsResult {
	var res []*AvgAggsResult
	for _, bucket := range a.Aggregations.Avgs.Buckets {
		res = append(res, &AvgAggsResult{
			Count:     bucket.DocCount,
			Value:     bucket.AvgValue.Value,
			StartTime: bucket.KeyAsString,
		})
	}
	return res
}
