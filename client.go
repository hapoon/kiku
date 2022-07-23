package kiku

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

const endpointURL = "https://atnd.ak4.jp/api/cooperation"

type client struct {
	c *http.Client
}

func (c client) Get(ctx context.Context, url string) (response *http.Response, err error) {
	reqURL := endpointURL + url

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return
	}

	log.Println(http.MethodGet, reqURL)
	response, err = c.c.Do(req)
	return
}

func (c client) Post(ctx context.Context, url string, body interface{}) (response *http.Response, err error) {
	reqURL := endpointURL + url

	b, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBuffer(b))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	log.Println(http.MethodPost, reqURL)
	response, err = c.c.Do(req)
	return
}

func (c client) Patch(ctx context.Context, url string, body interface{}) (response *http.Response, err error) {
	reqURL := endpointURL + url

	b, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, reqURL, bytes.NewBuffer(b))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	log.Println(http.MethodPatch, reqURL)
	response, err = c.c.Do(req)
	return
}

func (c client) Delete(ctx context.Context, url string, body interface{}) (response *http.Response, err error) {
	reqURL := endpointURL + url

	b, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, bytes.NewBuffer(b))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	log.Println(http.MethodDelete, reqURL)
	response, err = c.c.Do(req)
	return
}

func newClient() *client {
	return &client{
		c: &http.Client{},
	}
}
