package util

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Bundle struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Request struct {
	Url       string            `json:"url"`
	Method    string            `json:"method"`
	Headers   http.Header       `json:"headers"`
	PathItems map[string]string `json:"pathItems"`
	CreatedAt time.Time         `json:"createdAt"`
	request   *http.Request
}

type Response struct {
	Status          int         `json:"status"`
	Headers         http.Header `json:"headers"`
	ContentType     string      `json:"contentType"`
	Payload         any         `json:"payload"`
	ValidationError error       `json:"validationError"`
}

func NewRequest(req *http.Request) *Request {
	request := Request{request: req}
	request.Url = req.URL.String()
	request.Headers = req.Header
	request.CreatedAt = time.Now()
	return &request
}

func (r Request) Request() *http.Request {
	return r.request
}

func StatusCodeOrDefault(status string, def int) int {
	if val, err := strconv.Atoi(status); err == nil {
		return val
	}
	return def
}

func WriteJSON(ctx echo.Context, res *Response) error {
	for k, _ := range res.Headers {
		ctx.Response().Header().Set(k, res.Headers.Get(k))
	}
	if res.Payload != nil {
		if data, ok := res.Payload.([]byte); ok {
			return ctx.Blob(res.Status, res.ContentType, data)
		}
		return ctx.JSON(res.Status, res.Payload)
	} else {
		return ctx.NoContent(res.Status)
	}
}
