package utils

import (
	"context"
)

type HttpClient interface {
	Send(ctx context.Context, url, method, body string) ([]byte, error)
}
