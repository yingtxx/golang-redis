// Author: Steve Zhang
// Date: 2020/9/25 11:11 上午

package http

import "net/url"

func (ct *ClientContainer) Get(baseURL string, query url.Values, respst interface{}) (err error) {
	client := ct.MustGetClient()
	defer ct.PutClient(client)

	return client.Get(baseURL, query, respst)
}

func (ct *ClientContainer) Post(baseURL string, reqdata, respst interface{}) (err error) {
	client := ct.MustGetClient()
	defer ct.PutClient(client)

	return client.Post(baseURL, reqdata, respst)
}

func (ct *ClientContainer) SendRequest(method string, baseURL string, reqBody string, headerMap map[string]string, response interface{})(map[string]interface{}, error) {
	client := ct.MustGetClient()
	defer ct.PutClient(client)

	return client.SendRequest(method, baseURL, reqBody, headerMap, response)
}
