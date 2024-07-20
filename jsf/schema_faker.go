package jsf

import (
	"github.com/brianvoe/gofakeit"
	"github.com/getkin/kin-openapi/openapi3"
	regen "github.com/zach-klippenstein/goregen"
	"math/rand"
	"mockambo/evaluator"
	"mockambo/extension"
	"mockambo/util"
	"slices"
)

const RFC3339local = "2006-01-02T15:04:05Z"

// GenerateDataFromSchema recursively generates data from an OpenAPI schema
func GenerateDataFromSchema(schema *openapi3.Schema, nMext extension.Mext, ev evaluator.Evaluator) (any, error) {
	if schema == nil {
		schema = &openapi3.Schema{}
	}
	var mext extension.Mext
	if schema.Extensions != nil {
		mx, err := extension.MergeMextWithExtensions(nMext, schema.Extensions)
		if err != nil {
			return nil, err
		}
		mext = mx
	} else {
		mext = nMext
	}

	return generateByPriority(schema, mext, ev)
}

// generateString generates a random string, respecting the requirements of the schema
func generateString(s *openapi3.Schema, mext extension.Mext) (string, error) {
	var err error
	res := ""
	ln := int(s.MinLength)
	mx := 16
	if ln > mx {
		mx = ln
	}
	if s.MaxLength != nil {
		mx = int(*s.MaxLength)
	}
	ln = mx
	if len(s.Pattern) > 0 {
		if res, err = regen.Generate(s.Pattern); err != nil {
			return res, err
		}
	} else if len(s.Format) > 0 {
		switch s.Format {
		case "date-time":
			return gofakeit.Date().Format(RFC3339local), nil
		case "uri-template", "uri":
			return gofakeit.URL(), nil
		}
	} else {
		for i := 0; len(res) < ln; i++ {
			if i > 0 {
				res += " "
			}
			res += gofakeit.Word()
		}
	}
	if len(res) > ln {
		res = res[0:ln]
	}
	return res, err
}

// generateInt generates a random integer, respecting the requirements of the schema
func generateInt(schema *openapi3.Schema, mext extension.Mext) int {
	mn := 0
	mx := 100
	if schema.Min != nil {
		mn = int(*schema.Min)
	}
	if schema.Max != nil {
		mx = int(*schema.Max)
	}
	return rand.Intn(mx-mn) + mn
}

// generateFloat generates a random float, respecting the requirements of the schema
func generateFloat(schema *openapi3.Schema, mext extension.Mext) float64 {
	var mn float64 = 0
	var mx float64 = 100
	if schema.Min != nil {
		mn = *schema.Min
	}
	if schema.Max != nil {
		mx = *schema.Max
	}
	v := rand.Float64()
	if v < mn {
		v += mn
	}
	if v > mx {
		v = mx
	}
	return v
}

// generateByPriority generates the right type of data based on the schema requirements.
// Additionally, it will determine what's the correct generation methodology using the PayloadGenerationModes
// priority system
func generateByPriority(schema *openapi3.Schema, mext extension.Mext, ev evaluator.Evaluator) (any, error) {
	if schema.Type == nil {
		// rarely we can end up with a a schema == nil. We change it to TypeObject because it's the safest
		schema.Type = &openapi3.Types{openapi3.TypeObject}
	}
	if schema.Enum != nil {
		return schema.Enum[rand.Intn(len(schema.Enum))], nil
	}
	for _, m := range mext.PayloadGenerationModes {
		switch m {
		case extension.ModeDefault:
			if schema.Default != nil {
				return schema.Default, nil
			}
		case extension.ModeExample:
			if schema.Example != nil {
				return schema.Example, nil
			}
		case extension.ModeFaker:
			if mext.Faker != "" {
				return Fake(mext.Faker), nil
			}
		case extension.ModeTemplate:
			if mext.Template != "" {
				return ev.Template(mext.Template)
			}
		case extension.ModeSchema:
			if schema.Type.Includes(openapi3.TypeString) {
				return generateString(schema, mext)
			}
			if schema.Type.Includes(openapi3.TypeInteger) {
				return generateInt(schema, mext), nil
			}
			if schema.Type.Includes(openapi3.TypeNumber) {
				if schema.Format == "float" {
					return generateFloat(schema, mext), nil
				}
				return generateInt(schema, mext), nil
			}
			if schema.Type.Includes(openapi3.TypeBoolean) {
				return gofakeit.Bool(), nil
			}

			if schema.OneOf != nil {
				return GenerateDataFromSchema(schema.OneOf[rand.Intn(len(schema.OneOf))].Value, mext, ev)
			}
			if schema.AnyOf != nil {
				return GenerateDataFromSchema(schema.OneOf[rand.Intn(len(schema.AnyOf))].Value, mext, ev)
			}
			if schema.AllOf != nil {
				sx := &openapi3.Schema{}
				MergeAllOf(schema.AllOf, 0, sx)
				return GenerateDataFromSchema(sx, mext, ev)
			}

			if schema.Type.Includes(openapi3.TypeObject) {
				res := make(map[string]any)
				for k, p := range schema.Properties {
					mx, err := extension.MergeMextWithExtensions(mext, p.Value.Extensions)
					if err != nil {
						return res, err
					}
					if mx.Display || util.RequiredOrRandom(slices.Contains(schema.Required, k)) {
						var err error
						if res[k], err = GenerateDataFromSchema(p.Value, mx, ev); err != nil {
							return res, err
						}
					}
				}
				if schema.AdditionalProperties.Schema != nil {
					// if AdditionalProperties has a schema, it means it wants to represent a generic map
					// therefore we generate a key/value to make it happy
					var err error
					if res[gofakeit.HipsterWord()], err = GenerateDataFromSchema(schema.AdditionalProperties.Schema.Value, mext, ev); err != nil {
						return res, err
					}
				}
				return res, nil
			}
			if schema.Type.Includes(openapi3.TypeArray) {
				res := make([]any, 0)
				for range ComputeItemsCount(schema) {
					item, err := GenerateDataFromSchema(schema.Items.Value, mext, ev)
					res = append(res, item)
					if err != nil {
						return res, err
					}
				}
				return res, nil
			}
		case extension.ModeScript:
			if mext.Script != "" {
				return ev.RunScript(mext.Script)
			}
		}
	}
	return nil, nil
}
