package util

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"net/http"
)

func HttpClient() *resty.Client {
	client := resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})

	return client
}
