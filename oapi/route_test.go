package oapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"mockambo/util"
	"net/http"
	"os"
	"testing"
)

func TestRouteDef_SelectResponse(t *testing.T) {
	data, _ := os.ReadFile("../test_data/petstore.yaml")
	doc, _ := NewDoc(data)
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
}

func TestRouteDef_Process(t *testing.T) {
	data, _ := os.ReadFile("../test_data/petstore.yaml")
	doc, _ := NewDoc(data)
	req, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	request := util.NewRequest(req)
	route, _ := doc.FindRoute(request)
	out, _ := route.Process(context.Background(), request)
	assert.NotNil(t, out.Payload)
}
