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
