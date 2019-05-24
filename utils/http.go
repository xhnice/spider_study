package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync_study/logger"
)

// HTTPRequest -
func HTTPRequest(method, uri string, header map[string]string, data []byte) ([]byte, error) {
	resp, err := HTTPResponse(method, uri, header, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("读取接口[%s][%s]返回的结果失败~%s", method, uri, err)
		return nil, err
	}

	return body, nil
}

// HTTPResponse -
func HTTPResponse(method, uri string, header map[string]string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, uri, strings.NewReader(string(data)))
	if err != nil {
		logger.Errorf("构建请求[%s][%s]失败~%s", method, uri, err)
		return nil, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("请求接口[%s][%s]失败~%s", method, uri, err)
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		logger.Errorf("请求接口[%s][%s]失败~未找到页面", method, uri)
		return nil, errors.New("page not found")
	}
	return resp, nil
}

// HTTPGet -
func HTTPGet(uri string, params url.Values) ([]byte, error) {
	if params != nil {
		urlParam := params.Encode()
		if strings.Contains(uri, "?") {
			uri += "&" + urlParam
		} else {
			uri += "?" + urlParam
		}
	}
	return HTTPRequest(http.MethodGet, uri, nil, nil)
}

// HTTPGetResponse -
func HTTPGetResponse(uri string, params url.Values) (*http.Response, error) {
	if params != nil {
		urlParam := params.Encode()
		if strings.Contains(uri, "?") {
			uri += "&" + urlParam
		} else {
			uri += "?" + urlParam
		}
	}
	return HTTPResponse(http.MethodGet, uri, nil, nil)
}
