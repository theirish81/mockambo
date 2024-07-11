package oapi

import (
	"github.com/stretchr/testify/assert"
	"mockambo/util"
	"net/http"
	"os"
	"testing"
)

func TestDoc_FindRoute(t *testing.T) {
	data, _ := os.ReadFile("../test_data/petstore.yaml")
	doc, _ := NewDoc(data)
	req, _ := http.NewRequest("GET", "http://example.com/api/v3/pet/123", nil)
	route, err := doc.FindRoute(util.NewRequest(req))
	assert.Nil(t, err)
	assert.Equal(t, "getPetById", route.OperationID())
}
