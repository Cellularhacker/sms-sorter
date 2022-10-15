package proxy2

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

var DefaultClient *http.Client

func init() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100
	t.IdleConnTimeout = 90 * time.Second
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	DefaultClient = &http.Client{
		Timeout:   60 * time.Second,
		Transport: t,
	}
}

func Get(uri string) ([]byte, error) {
	resp, err := DefaultClient.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func Post(uri string, contentType string, body io.Reader) ([]byte, error) {
	resp, err := DefaultClient.Post(uri, contentType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
