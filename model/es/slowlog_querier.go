package es

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/dbrainhub/dbrainhub/api"
	"github.com/dbrainhub/dbrainhub/errors"
	esClient "github.com/elastic/go-elasticsearch/v8"
	"github.com/golang/protobuf/jsonpb"
)

type (
	SlowlogQuerier interface {
		Query(ctx context.Context, param *SlowlogQuerierParam) (*api.SearchMemberLogCountResponse, error)
	}

	SlowlogQuerierParam struct {
		Index   string
		Begin   string
		End     string
		Buckets int64
		Size    int64
		From    int64

		Cond map[string]interface{}
	}
)

func NewSlowlogQuerier(client *esClient.Client) SlowlogQuerier {
	return &slowlogQuerier{
		client: client,
	}
}

type slowlogQuerier struct {
	client *esClient.Client
}

func (s *slowlogQuerier) Query(ctx context.Context, param *SlowlogQuerierParam) (*api.SearchMemberLogCountResponse, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"from":             param.From,
		"size":             param.Size,
		"track_total_hits": true,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []interface{}{
					// bool-filter
					map[string]interface{}{"bool": map[string]interface{}{
						"filter": genMatches(param.Cond),
					}},

					// range
					map[string]interface{}{
						"range": genRange(param.Begin, param.End),
					},
				},
			},
		},
		"aggs": map[string]interface{}{
			"logs": map[string]interface{}{
				"date_histogram": genHistogram(param.Begin, param.End, param.Buckets),
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(param.Index),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			return nil, errors.FilebeatSearchError(e["error"].(map[string]interface{})["reason"].(string))
		}
	}

	searchRes := new(api.SearchMemberLogCountResponse)
	unm := jsonpb.Unmarshaler{
		AllowUnknownFields: true,
	}
	if err := unm.Unmarshal(res.Body, searchRes); err != nil {
		return nil, err
	}
	return searchRes, nil
}
