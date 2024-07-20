package jsf

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMergeAllOf(t *testing.T) {
	allOfs := openapi3.SchemaRefs{
		{
			Value: &openapi3.Schema{
				Type:     &openapi3.Types{openapi3.TypeObject},
				Required: []string{"foo"},
				Properties: openapi3.Schemas{
					"foo": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: &openapi3.Types{openapi3.TypeString},
						},
					},
				},
			},
		},
		{
			Value: &openapi3.Schema{
				Type:     &openapi3.Types{openapi3.TypeObject},
				Required: []string{"bar"},
				Properties: openapi3.Schemas{
					"bar": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: &openapi3.Types{openapi3.TypeString},
						},
					},
				},
			},
		},
	}
	out := &openapi3.Schema{}
	MergeAllOf(allOfs, 0, out)
	assert.NotNil(t, out.Properties["foo"])
	assert.NotNil(t, out.Properties["bar"])
}
