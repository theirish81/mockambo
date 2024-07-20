package oapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"mockambo/util"
	"net/http"
	"testing"
)

func TestResponseDef_GenerateResponseBundle(t *testing.T) {
	doc, _ := NewDoc("../test_data/petstore.yaml")
	r, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	req := util.NewRequest(r)
	route, _ := doc.FindRoute(req)
	out, _ := route.Process(context.Background(), req)
	assert.IsType(t, map[string]any{}, out.Payload)
	name, _ := out.Payload.(map[string]any)["name"]
	assert.Greater(t, len(name.(string)), 0)
	photoUrls, _ := out.Payload.(map[string]any)["photoUrls"]
	assert.IsType(t, []any{}, photoUrls)
}

func TestResponseDef_MediaExample(t *testing.T) {
	doc, _ := NewDoc("../test_data/custom.yaml")
	r, _ := http.NewRequest("GET", "http://localhost/api/v3/media-example1", nil)
	req := util.NewRequest(r)
	route, _ := doc.FindRoute(req)
	res, _ := route.selectResponse()
	out, _ := res.GenerateResponse(route.mext)
	assert.IsType(t, map[string]any{}, out.Payload)
	conv := out.Payload.(map[string]any)
	assert.Equal(t, "bar", conv["foo"])

	r, _ = http.NewRequest("GET", "http://localhost/api/v3/media-example2", nil)
	req = util.NewRequest(r)
	route, _ = doc.FindRoute(req)
	res, _ = route.selectResponse()
	out, _ = res.GenerateResponse(route.mext)
	assert.IsType(t, map[string]any{}, out.Payload)
	conv = out.Payload.(map[string]any)
	assert.Equal(t, "bar", conv["foo"])

	r, _ = http.NewRequest("GET", "http://localhost/api/v3/media-example2", nil)
	req = util.NewRequest(r)
	route, _ = doc.FindRoute(req)
	res, _ = route.selectResponse()
	res.r.Content["application/json"].Extensions["x-mockambo"].(map[string]any)["mediaExampleSelectorScript"] = ""
	out, _ = res.GenerateResponse(route.mext)
	assert.IsType(t, map[string]any{}, out.Payload)
	conv = out.Payload.(map[string]any)
	assert.Equal(t, "bar", conv["foo"])
}
