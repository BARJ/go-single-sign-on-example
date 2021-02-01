package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var DefaultHttpClient = &HttpClient{client: http.Client{Timeout: 10 * time.Second}}

type HttpClient struct {
	client http.Client
}

func (h HttpClient) HttpRequestJson(
	method string,
	url string,
	headers map[string]string,
	body string,
	v interface{},
) error {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	res, err := h.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return errors.New(res.Status)
	}

	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, v)
	if err != nil {
		return err
	}
	return nil
}

func HttpReplyError(w http.ResponseWriter, statusCode int, err error) {
	text := http.StatusText(statusCode)
	if err != nil {
		text = fmt.Sprintf("%s: %v", text, err)
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(text))
}

func HttpReplyJson(w http.ResponseWriter, statusCode int, media interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(media)
	if err != nil {
		HttpReplyError(w, http.StatusInternalServerError, nil)
		return
	}

	w.WriteHeader(statusCode)
	w.Write(data)
}
