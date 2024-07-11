package util

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type rw struct {
	data   []byte
	status int
}

func (w *rw) Header() http.Header {
	return http.Header{}
}
func (w *rw) Write(data []byte) (int, error) {
	w.data = data
	return len(data), nil
}
func (w *rw) WriteHeader(statusCode int) {
	w.status = statusCode
}

func TestWriteJSON(t *testing.T) {
	e := echo.New()
	u, _ := http.NewRequest("GET", "http://example.com", nil)
	x := &rw{}
	ctx := e.NewContext(u, x)
	res := Response{
		Payload: map[string]string{"foo": "bar"},
		Headers: http.Header{},
	}
	_ = WriteJSON(ctx, &res)
	assert.Equal(t, "{\"foo\":\"bar\"}\n", string(x.data))
	x.data = []byte{}
	res.Payload = nil
	_ = WriteJSON(ctx, &res)
	assert.Equal(t, "", string(x.data))
}

func TestStatusCodeOrDefault(t *testing.T) {
	assert.Equal(t, 400, StatusCodeOrDefault("400", 200))
	assert.Equal(t, 200, StatusCodeOrDefault("", 200))
}
