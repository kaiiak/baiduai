package baiduai

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	defaultClient     = &http.Client{Timeout: 90 * time.Second}
	defaultRetryCount = 3
)

func httpGet(url string, header map[string]string, query url.Values) (Scaner, error) {
	var (
		resp *http.Response
		req  *http.Request
		err  error
		b    []byte
	)
	defer logrus.Debugf("httpGet url [%s] header [%v] query [%v], at %s", url, header, query, time.Now())
	req, err = wrapRequest(http.MethodGet, url, header, query, nil)
	if err != nil {
		return nil, err
	}
	for i := 0; i < defaultRetryCount; i++ {
		resp, err = defaultClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}
	}
	if err != nil {
		logrus.Errorf("http get url [%s] error [%v]", req.URL.String(), err)
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("httpGet read [%s] response body error %v", req.URL.String(), err)
		return nil, err
	}

	return bytesScaner(b), nil
}

func wrapRequest(method, url string, header map[string]string, query url.Values, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for c, v := range header {
		req.Header.Set(c, v)
	}
	q := req.URL.Query()
	for k, v := range query {
		q[k] = v
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

// httpPost post not retry
func httpPost(url string, header map[string]string, query url.Values, body io.Reader) (Scaner, error) {
	req, err := wrapRequest(http.MethodPost, url, header, query, body)
	if err != nil {
		return nil, err
	}
	resp, err := defaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytesScaner(b), nil
}
