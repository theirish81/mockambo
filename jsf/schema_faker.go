package jsf

import (
	"github.com/brianvoe/gofakeit"
	"github.com/dop251/goja"
	"github.com/getkin/kin-openapi/openapi3"
	regen "github.com/zach-klippenstein/goregen"
	"math/rand"
	"mockambo/extension"
	"mockambo/util"
	"slices"
)

const RFC3339local = "2006-01-02T15:04:05Z"

func GenerateDataFromSchema(schema *openapi3.Schema, defaultMext extension.Mext, vm *goja.Runtime) (any, error) {
	mext, err := extension.MergeDefaultMextWithExtensions(defaultMext, schema.Extensions)
	if err != nil {
		return nil, err
	}
	return generateByPriority(schema, mext, vm)
}

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
		case "uri-template":
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

func generateByPriority(schema *openapi3.Schema, mext extension.Mext, vm *goja.Runtime) (any, error) {
	if schema.Enum != nil {
		return schema.Enum[rand.Intn(len(schema.Enum))], nil
	}
	for _, m := range mext.PayloadGenerationModes {
		switch m {
		case "default":
			if schema.Default != nil {
				return schema.Default, nil
			}
		case "example":
			if schema.Example != nil {
				return schema.Example, nil
			}
		case "schema":
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
			if schema.Type.Includes(openapi3.TypeObject) {
				res := make(map[string]any)
				for k, p := range schema.Properties {
					mx, err := extension.MergeDefaultMextWithExtensions(mext, p.Value.Extensions)
					if err != nil {
						return res, err
					}
					if mx.Display || util.RequiredOrRandom(slices.Contains(schema.Required, k)) {
						var err error
						if res[k], err = GenerateDataFromSchema(p.Value, mx, vm); err != nil {
							return res, err
						}
					}
				}
				return res, nil
			}
			if schema.Type.Includes(openapi3.TypeArray) {
				res := make([]any, 0)
				item, err := GenerateDataFromSchema(schema.Items.Value, mext, vm)
				res = append(res, item)
				return res, err
			}
		case "script":
			if mext.Script != nil {
				return vm.RunString(*mext.Script)
			}
		}
	}
	return nil, nil
}
