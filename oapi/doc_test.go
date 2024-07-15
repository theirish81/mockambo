package oapi

import (
	"github.com/stretchr/testify/assert"
	"mockambo/util"
	"net/http"
	"testing"
)

func TestDoc_FindRoute(t *testing.T) {
	doc, _ := NewDoc("../test_data/petstore.yaml")
	req, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	route, err := doc.FindRoute(util.NewRequest(req))
	assert.Nil(t, err)
	assert.Equal(t, "getPetById", route.OperationID())
}

func TestDoc_Servers(t *testing.T) {
	doc, _ := NewDoc("../test_data/petstore.yaml")
	servers := doc.Servers()
	assert.Len(t, servers, 1)
}
