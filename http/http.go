package http

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"
)

func Get(url string, headers http.Header) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr, Timeout: time.Second * 30}
	defer client.CloseIdleConnections()
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		if len(request.Header.Get(k)) > 0 {
			request.Header.Set(k, strings.Join(v, ","))
		} else {
			request.Header.Add(k, strings.Join(v, ","))
		}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	return respBody, err
}

func Post(url string, headers http.Header, body io.Reader) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: tr, Timeout: time.Second * 30}
	defer client.CloseIdleConnections()
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		if len(request.Header.Get(k)) > 0 {
			request.Header.Set(k, strings.Join(v, ","))
		} else {
			request.Header.Add(k, strings.Join(v, ","))
		}
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	return respBody, err
}
