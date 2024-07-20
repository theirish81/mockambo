package jsf

import (
	"github.com/getkin/kin-openapi/openapi3"
	"math/rand"
)

func ComputeItemsCount(schema *openapi3.Schema) int {
	mn := int(schema.MinItems)
	mx := 1
	if schema.MaxItems != nil {
		mx = int(*schema.MaxItems)
	}
	count := mn
	if mx > count {
		count = rand.Intn(mx-mn) + mn
	}
	if count == 0 {
		count = 1
	}
	return count
}

// MergeAllOf will merge into one schemas all the sub-schemas declared in the AllOf
func MergeAllOf(allOf openapi3.SchemaRefs, index int, schema *openapi3.Schema) {
	extract := allOf[index].Value
	if extract.Type.Includes(openapi3.TypeObject) {
		if schema.Properties == nil {
			schema.Type = &openapi3.Types{openapi3.TypeObject}
			schema.Properties = openapi3.Schemas{}
		}
		if schema.Required == nil {
			schema.Required = make([]string, 0)
		}
		for k, v := range extract.Properties {
			schema.Properties[k] = v
		}
		if extract.Required != nil {
			schema.Required = append(schema.Required, extract.Required...)
		}
	}
	if index < len(allOf)-1 {
		MergeAllOf(allOf, index+1, schema)
	}
}
