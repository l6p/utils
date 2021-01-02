package json

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
)

type Request interface {
	D(data *Data) Request
	J(json string) Request
	Header(key string, value string) Request
	Ignore() Request
	Get(url string) Response
	Post(url string) Response
	Put(url string) Response
	Patch(url string) Response
	Delete(url string) Response
	Do(method Method, url string) Response
	Redo() Response
}

type RequestImpl struct {
	client  *Client
	data    *Data
	headers map[string]string
	method  Method
	url     string
}

func (r *RequestImpl) D(data *Data) Request {
	r.data = data
	return r
}

func (r *RequestImpl) J(json string) Request {
	return r.D(D(json))
}

func (r *RequestImpl) Header(key string, value string) Request {
	r.headers[key] = value
	return r
}

func (r *RequestImpl) Ignore() Request {
	return r.Header("x-l6p-ignore", "true")
}

func (r *RequestImpl) Get(url string) Response {
	return r.Do(GET, url)
}

func (r *RequestImpl) Post(url string) Response {
	return r.Do(POST, url)
}

func (r *RequestImpl) Put(url string) Response {
	return r.Do(PUT, url)
}

func (r *RequestImpl) Patch(url string) Response {
	return r.Do(PATCH, url)
}

func (r *RequestImpl) Delete(url string) Response {
	return r.Do(DELETE, url)
}

func (r *RequestImpl) Do(method Method, url string) Response {
	var jsonObj interface{}
	var resp *resty.Response
	var err error

	r.method = method
	r.url = url

	req := r.client.rc.R().
		SetHeaders(r.headers).
		SetResult(&jsonObj)

	switch method {
	case GET:
		resp, err = req.Get(url)
	case POST:
		resp, err = req.SetBody(r.data.GetJson("")).Post(url)
	case PUT:
		resp, err = req.SetBody(r.data.GetJson("")).Put(url)
	case PATCH:
		resp, err = req.SetBody(r.data.GetJson("")).Patch(url)
	case DELETE:
		resp, err = req.Patch(url)
	default:
		panic(errors.New("invalid method"))
	}

	if err != nil {
		panic(err)
	}

	if resp.IsError() {
		panic(errors.New("status code >= 400"))
	}

	return &ResponseImpl{
		req:  r,
		data: &Data{jsonObj: &jsonObj},
	}
}

func (r *RequestImpl) Redo() Response {
	return r.Do(r.method, r.url)
}
