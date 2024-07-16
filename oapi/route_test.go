package oapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"mockambo/util"
	"net/http"
	"testing"
)

func TestRouteDef_SelectResponse(t *testing.T) {
	doc, _ := NewDoc("../test_data/petstore.yaml")
	req, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	request := util.NewRequest(req)
	route, _ := doc.FindRoute(request)
	_ = route.validateRequest(context.Background(), request)
	resp, _ := route.selectResponse()
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.status)

	req, _ = http.NewRequest("GET", "http://example.com/api/v3/pet/abc", nil)
	request = util.NewRequest(req)
	route, _ = doc.FindRoute(request)
	route.setValidationError(route.validateRequest(context.Background(), request))
	resp, _ = route.selectResponse()
	assert.NotNil(t, resp)
	assert.Equal(t, 400, resp.status)

	req, _ = http.NewRequest("GET", "http://example.com/api/v4/pet/abc", nil)
	request = util.NewRequest(req)
	route, err := doc.FindRoute(request)
	assert.NotNil(t, err)

	req, _ = http.NewRequest("POST", "http://example.com/api/v3/user", nil)
	request = util.NewRequest(req)
	route, err = doc.FindRoute(request)
	res, err := route.selectResponse()
	assert.Nil(t, err)
	assert.Equal(t, 200, res.status)
}

func TestRouteDef_Process(t *testing.T) {
	doc, _ := NewDoc("../test_data/petstore.yaml")
	req, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	request := util.NewRequest(req)
	route, _ := doc.FindRoute(request)
	out, _ := route.Process(context.Background(), request)
	assert.NotNil(t, out.Payload)
}
