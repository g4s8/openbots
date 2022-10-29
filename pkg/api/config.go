package api

import "time"

type Config struct {
	// Addr - HTTP server address.
	Addr string

	// ReadTimeout - HTTP server read timeout.
	ReadTimeout time.Duration

	// RequestTimeout - HTTP request timeout.
	RequestTimeout time.Duration
}
