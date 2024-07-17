package jsf

import "github.com/getkin/kin-openapi/openapi3"

func Merge(allOf openapi3.SchemaRefs, index int, schema *openapi3.Schema) {
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
		Merge(allOf, index+1, schema)
	}
}
