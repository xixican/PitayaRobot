package util

import (
	"github.com/topfreegames/pitaya/logger"
	"io"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string, authorization string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Log.Errorf("创建get request错误:%s, url=%s", err.Error(), url)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Errorf("调用get错误:%s, url=%s", err.Error(), url)
		return nil
	}
	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.Log.Errorf("调用get错误:%s, url=%s", err.Error(), url)
		return nil
	}
	return b
}

func HttpPost(url string, authorization string, body io.Reader) []byte {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		logger.Log.Errorf("创建post request错误:%s, url=%s", err.Error(), url)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Errorf("调用post错误:%s, url=%s", err.Error(), url)
		return nil
	}
	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.Log.Errorf("调用post错误:%s, url=%s", err.Error(), url)
		return nil
	}
	return b
}

func HttpPut(url string, authorization string, body io.Reader) []byte {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		logger.Log.Errorf("创建put request错误:%s, url=%s", err.Error(), url)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Log.Errorf("调用put错误:%s, url=%s", err.Error(), url)
		return nil
	}
	b, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logger.Log.Errorf("调用put错误:%s, url=%s", err.Error(), url)
		return nil
	}
	return b
}
