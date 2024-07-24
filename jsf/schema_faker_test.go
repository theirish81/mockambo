package jsf

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"mockambo/evaluator"
	"mockambo/extension"
	"testing"
	"time"
)

func TestGenerateDataFromSchema(t *testing.T) {
	doc, _ := openapi3.NewLoader().LoadFromFile("../test_data/petstore.yaml")
	path := doc.Paths.Value("/pet/{petId}")
	mext, _ := extension.NewMextFromExtensions(nil)
	ev := evaluator.NewEvaluator()
	ev.Set("fake", Fake)
	ev.Set("pathItems", map[string]any{"petId": 123})
	out, _ := GenerateDataFromSchema(path.Get.Responses.Value("200").Value.Content.Get("application/json").Schema.Value, mext, ev)
	assert.IsType(t, map[string]any{}, out)
	name, _ := out.(map[string]any)["name"]
	assert.Greater(t, len(name.(string)), 0)
	photoUrls, _ := out.(map[string]any)["photoUrls"]
	assert.IsType(t, []any{}, photoUrls)
}

func TestGenerateString(t *testing.T) {
	var maxLen uint64 = 8
	res, _ := generateString(&openapi3.Schema{
		Type:      &openapi3.Types{"string"},
		MinLength: 3,
		MaxLength: &maxLen,
	}, extension.Mext{})
	assert.GreaterOrEqual(t, len(res), 3)
	assert.LessOrEqual(t, len(res), int(maxLen))

	res, _ = generateString(&openapi3.Schema{
		Type:   &openapi3.Types{"string"},
		Format: "date-time",
	}, extension.Mext{})
	_, err := time.Parse(RFC3339local, res)
	assert.Nil(t, err)
}

func TestGenerateFloat(t *testing.T) {
	res := generateFloat(&openapi3.Schema{
		Type:   &openapi3.Types{"number"},
		Format: "float",
	}, extension.Mext{})
	assert.IsType(t, 1.5, res)

	mn := 0.5
	mx := 1.2
	res = generateFloat(&openapi3.Schema{
		Type:   &openapi3.Types{"number"},
		Format: "float",
		Min:    &mn,
		Max:    &mx,
	}, extension.Mext{})
	assert.GreaterOrEqual(t, res, 0.5)
	assert.LessOrEqual(t, res, 1.2)

	res = generateFloat(&openapi3.Schema{
		Type:   &openapi3.Types{"number"},
		Format: "foo",
	}, extension.Mext{})
	assert.IsType(t, 1.2, res)
}

func TestGenerateInt(t *testing.T) {
	res := generateInt(&openapi3.Schema{
		Type: &openapi3.Types{"integer"},
	}, extension.Mext{})
	assert.IsType(t, 1, res)

	var mn float64 = 1
	var mx float64 = 5
	res = generateInt(&openapi3.Schema{
		Type: &openapi3.Types{"integer"},
		Min:  &mn,
		Max:  &mx,
	}, extension.Mext{})
	assert.GreaterOrEqual(t, res, 1)
	assert.LessOrEqual(t, res, 5)
}

func TestAdditionalProperties(t *testing.T) {
	l := openapi3.NewLoader()
	l.IsExternalRefsAllowed = true
	doc, err := l.LoadFromFile("../test_data/custom.yaml")
	fmt.Println(err)
	path := doc.Paths.Value("/additional-properties")
	mext, _ := extension.NewMextFromExtensions(nil)
	ev := evaluator.NewEvaluator()
	out, _ := GenerateDataFromSchema(path.Get.Responses.Value("200").Value.Content.Get("application/json").Schema.Value, mext, ev)
	assert.IsType(t, map[string]any{}, out)
	for k, v := range out.(map[string]any) {
		assert.IsType(t, "foo", k)
		assert.IsType(t, 22, v)
	}
}

func TestOneOf(t *testing.T) {
	oneOf := openapi3.SchemaRefs{
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
	mext, _ := extension.NewMextFromExtensions(nil)
	s, err := GenerateDataFromSchema(&openapi3.Schema{
		OneOf: oneOf,
	}, mext, evaluator.NewEvaluator())
	assert.Nil(t, err)
	assert.IsType(t, map[string]any{}, s)
	_, ok1 := s.(map[string]any)["foo"]
	_, ok2 := s.(map[string]any)["bar"]
	assert.True(t, ok1 || ok2)
	assert.False(t, ok1 && ok2)
}

func TestAllOf(t *testing.T) {
	allOf := openapi3.SchemaRefs{
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
	mext, _ := extension.NewMextFromExtensions(nil)
	s, err := GenerateDataFromSchema(&openapi3.Schema{
		AllOf: allOf,
	}, mext, evaluator.NewEvaluator())
	assert.Nil(t, err)
	assert.IsType(t, map[string]any{}, s)
	_, ok1 := s.(map[string]any)["foo"]
	_, ok2 := s.(map[string]any)["bar"]
	assert.True(t, ok1 && ok2)
}
