package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var (
	ContentTypeTextXml = "text/xml"
	ContentTypeHtml    = "text/html; charset=utf-8"
	ContentTypeTextCss = "text/css; charset=utf-8"
	ContentTypeXJS     = "application/x-javascript"
	ContentTypeJS      = "text/javascript"
	ContentTypeJson    = "application/json; charset=utf-8"
	ContentTypeForm    = "application/x-www-form-urlencoded"
	ContentTypeImg     = "image/png"
)

// PostJSON - send an http post json Request.
func PostJSON(url string, data interface{}, deadline, dialTimeout time.Duration) ([]byte, int, error) {
	buf, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	return Request(http.MethodPost, url, bytes.NewBuffer(buf), deadline, dialTimeout, map[string]string{"Content-Type": ContentTypeJson})
}

// PostForm - sen an http post form Request
func PostForm(url string, data []byte, deadline, dialTimeout time.Duration) ([]byte, int, error) {
	return Request(http.MethodPost, url, bytes.NewBuffer(data), deadline, dialTimeout, map[string]string{"Content-Type": ContentTypeForm})
}

// SimpleGet - send an http get Request
func SimpleGet(url string, deadline, dialTimeout time.Duration) ([]byte, int, error) {
	return Request(http.MethodGet, url, bytes.NewBuffer(nil), deadline, dialTimeout, nil)
}

// SimpleDelete - send an simple http delete Request
func SimpleDelete(url string, deadline, dialTimeout time.Duration) ([]byte, int, error) {
	return Request(http.MethodDelete, url, bytes.NewBuffer(nil), deadline, dialTimeout, nil)
}

// SimplePut - send an simple http put Request
func SimplePut(url string, deadline, dialTimeout time.Duration) ([]byte, int, error) {
	return Request(http.MethodPut, url, bytes.NewBuffer(nil), deadline, dialTimeout, nil)
}

// Request - send an http Request
func Request(method, url string, body io.Reader, deadline, dialTimeout time.Duration, header map[string]string) ([]byte, int, error) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(deadline)
				c, err := net.DialTimeout(netw, addr, dialTimeout)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, 0, err
	}

	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return data, resp.StatusCode, nil
}

// HttpResponse - htt Response
type HttpResponse struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data,omitempty"`
}

// NewHttpResponse -
func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		Code:    0,
		Message: "success",
	}
}

// Response - write data to resp
func (h *HttpResponse) Response(resp http.ResponseWriter) {
	resp.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(h)
	resp.Write(data)
}

// ResponseWithErr - write data to resp with error
func (h *HttpResponse) ResponseWithErr(resp http.ResponseWriter, err error) {
	resp.WriteHeader(http.StatusOK)
	if err != nil {
		h.Error(err)
	}

	data, _ := json.Marshal(h)
	resp.Write(data)
}

// Error - set Error
func (h *HttpResponse) Error(err error) {
	h.Code = 1
	h.Message = err.Error()
}
