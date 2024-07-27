package oapi

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"mockambo/util"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestRouteDef_SelectResponse(t *testing.T) {
	doc, _ := NewDoc("../test_data/petstore.yaml", "")
	req, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	request := util.NewRequest(req)
	route, _ := doc.FindRoute(request)
	route.initRequestValidationInput(request)
	_ = route.validateRequest(context.Background())
	resp, _ := route.selectResponse()
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.status)

	req, _ = http.NewRequest("GET", "http://example.com/api/v3/pet/abc", nil)
	request = util.NewRequest(req)
	route, _ = doc.FindRoute(request)
	route.initRequestValidationInput(request)
	route.setValidationError(route.validateRequest(context.Background()))
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

	req, _ = http.NewRequest("DELETE", "http://example.com/api/v3/pet/abc", nil)
	request = util.NewRequest(req)
	route, err = doc.FindRoute(request)
	res, err = route.selectResponse()
	assert.Nil(t, err)
	assert.Equal(t, 400, res.status)
}

func TestRouteDef_Process(t *testing.T) {
	doc, _ := NewDoc("../test_data/petstore.yaml", "")
	req, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	request := util.NewRequest(req)
	route, _ := doc.FindRoute(request)
	out, _ := route.Process(context.Background(), request)
	assert.NotNil(t, out.Payload)
}

func TestRouteDef_RecordPlayback(t *testing.T) {
	doc, _ := NewDoc("../test_data/github.yaml", "")
	doc.mext.RecordingPath = path.Join(os.TempDir(), "mockambo_"+gofakeit.UUID())
	req, _ := http.NewRequest("GET", "http://localhost:8080/orgs/github/repos", nil)
	request := util.NewRequest(req)
	route, _ := doc.FindRoute(request)
	out, _ := route.Process(context.Background(), request)
	assert.NotNil(t, out.Payload)

	// regenerating request as it may have been altered by 302. This can't happen in the real world because
	// requests do not get reused
	req, _ = http.NewRequest("GET", "http://localhost:8080/orgs/github/repos", nil)
	request = util.NewRequest(req)
	out, _ = route.Process(context.Background(), request)
	assert.NotNil(t, out.Payload)
	assert.NotEmpty(t, out.Headers.Get("x-mockambo-playback"))

	req, _ = http.NewRequest("GET", "http://localhost:8080/orgs/github/repos", nil)
	req.Header.Set("x-mockambo-invalidate-recording", "true")
	request = util.NewRequest(req)
	out, _ = route.Process(context.Background(), request)
	assert.NotNil(t, out.Payload)
	assert.Empty(t, out.Headers.Get("x-mockambo-playback"))

}
