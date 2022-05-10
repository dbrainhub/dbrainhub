package es

import (
	"github.com/dbrainhub/dbrainhub/utils/search_time"
)

// gen muiti terms statement
func genTerms(cond map[string]interface{}) []map[string]interface{} {
	res := []map[string]interface{}{}
	for key, value := range cond {
		res = append(res, map[string]interface{}{
			"term": map[string]interface{}{
				key: value,
			},
		})
	}
	return res
}

// gen muiti matches statement
func genMatches(cond map[string]interface{}) []map[string]interface{} {
	res := []map[string]interface{}{}
	for key, value := range cond {
		res = append(res, map[string]interface{}{
			"match": map[string]interface{}{
				key: value,
			},
		})
	}
	return res
}

// from/to format "2006-01-02T15:04:05.999Z"
func genRange(from, to string) map[string]interface{} {
	return map[string]interface{}{
		"@timestamp": map[string]interface{}{
			"format": "strict_date_optional_time",
			"gte":    from,
			"lte":    to,
		},
	}
}

func genHistogram(from, to string, buckets int64) map[string]interface{} {
	return map[string]interface{}{
		"date_histogram": map[string]interface{}{
			"field":          "@timestamp",
			"fixed_interval": search_time.GetInterval(from, to, buckets),
			"min_doc_count":  1, // ignore count == 0
			"format":         "yyyy-MM-dd HH:mm:ss",
			"time_zone":      "Asia/Shanghai",
		},
	}
}
