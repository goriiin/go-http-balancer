package domain

import "net/url"

type Backend struct {
	URL        *url.URL
	ConnsCount int32
	Index      int
}
