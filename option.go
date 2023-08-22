package ftp

import "time"

type Option struct {
	Filename string
	Addr     string
	Retries  int
	Timeout  time.Duration
}
