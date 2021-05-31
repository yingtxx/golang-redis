package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"yufu_server/log"
)

const (
	defaultContentType = "application/json"
)

type Client struct {
	*http.Client
}

type ClientConf struct {
	DialTimeoutSecond     int // 连接超时
	DialKeepAliveSecond   int // 开启长连接
	MaxIdleConns          int // 最大空闲连接数
	MaxIdleConnsPerHost   int // HOST最大空闲连接数
	IdleConnTimeoutSecond int // 空闲连接超时
}

func NewClient(cf *ClientConf) (cli *Client) {
	cli = &Client{
		Client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        cf.MaxIdleConns,
				MaxIdleConnsPerHost: cf.MaxIdleConnsPerHost,
				IdleConnTimeout:     time.Duration(cf.IdleConnTimeoutSecond) * time.Second,
				Proxy:               http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   time.Duration(cf.DialTimeoutSecond) * time.Second,
					KeepAlive: time.Duration(cf.DialKeepAliveSecond) * time.Second,
				}).DialContext,
			},
		},
	}
	return
}

func (cli *Client) Get(baseURL string, query url.Values, respst interface{}) (err error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		err = fmt.Errorf("url parse: %w", err)
		return
	}

	u.RawQuery = query.Encode()

	resp, err := cli.Client.Get(u.String())
	if err != nil {
		err = fmt.Errorf("client get: %w", err)
		return
	}

	if err = scanresp(resp, respst); err != nil {
		err = fmt.Errorf("scan response: %w", err)
		return
	}

	return
}

func (cli *Client) Post(baseURL string, reqdata, respst interface{}) (err error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		err = fmt.Errorf("url parse: %w", err)
		return
	}

	data, err := json.Marshal(reqdata)
	if err != nil {
		err = fmt.Errorf("json marshal: %w", err)
		return
	}

	resp, err := cli.Client.Post(u.String(), defaultContentType, bytes.NewReader(data))
	if err != nil {
		err = fmt.Errorf("client post: %w", err)
		return
	}

	if err = scanresp(resp, respst); err != nil {
		err = fmt.Errorf("scan response: %w", err)
		return
	}

	return
}

func (cli *Client) Close() (err error) {
	cli.Client.CloseIdleConnections()
	return
}

func scanresp(resp *http.Response, st interface{}) (err error) {
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("response code %d", resp.StatusCode)
		return
	}

	if resp.Body == nil {
		err = fmt.Errorf("body is nil")
		return
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(st); err != nil {
		err = fmt.Errorf("json decode: %w", err)
		return
	}

	return
}

func (cli *Client)SendRequest(method string, baseURL string, reqBody string, headerMap map[string]string, response interface{})(responseMap map[string]interface{}, err error){
	u, err := url.Parse(baseURL)
	if err != nil {
		err = fmt.Errorf("url parse: %w", err)
		return
	}

	//初始化请求
	req, httpErr := http.NewRequest(method, u.String(), strings.NewReader(reqBody))
	if httpErr != nil {
		err = fmt.Errorf("httpErr:" + httpErr.Error(), err)
		return
	}
	//设置请求头
	for headerKey:= range headerMap{
		req.Header.Set(headerKey, headerMap[headerKey])
	}

	resp, err := cli.Client.Do(req)
	if err != nil {
		err = fmt.Errorf("client "+ method +": %w", err)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	log.LoggerContainer.Info(map[string]interface{}{
		"request_uri" : "",
		"yufu_api" : req.URL.Host + req.URL.Path,
		"yufu_req" : reqBody,
		"yufu_res" : string(bodyBytes),
	})
	_ = json.Unmarshal(bodyBytes, &responseMap)
	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		err = fmt.Errorf("scan response: %w", err)
		return
	}

	return
}