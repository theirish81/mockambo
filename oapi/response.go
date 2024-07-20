package oapi

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"mockambo/evaluator"
	"mockambo/exceptions"
	"mockambo/extension"
	"mockambo/jsf"
	"mockambo/util"
	"net/http"
	"strings"
)

// ResponseDef is a response definition. It contains the openapi3 definition, plus its instrumentation
type ResponseDef struct {
	status      int
	r           *openapi3.Response
	err         error
	defaultMext extension.Mext
	evaluator   evaluator.Evaluator
}

func NewResponseDef(r *openapi3.Response, status int, mext extension.Mext, ev evaluator.Evaluator) (ResponseDef, error) {
	defaultMext, err := extension.MergeMextWithExtensions(mext, r.Extensions)
	ev.Set("status", status)
	return ResponseDef{r: r, status: status, defaultMext: defaultMext, evaluator: ev}, err
}

// determineJsonMediaType looks for a media definition that contains "json" in it. As Mockambo only support JSON
// this makes kinda sense. If none are found, then an empty string is returned
func (r ResponseDef) determineJsonMediaType() string {
	for k := range r.r.Content {
		if strings.Contains(k, "json") {
			return k
		}
	}
	return ""
}

// generateResponsePayload generates a sample payload for the response. Because there are multiple generation,
// strategies, some of which are scriptable, the response can be an object, a string  or a byte slice. Hence
// the `any` return type
func (r ResponseDef) generateResponsePayload(mext extension.Mext) (any, error) {
	if mediaType := r.determineJsonMediaType(); mediaType != "" {
		mext, err := extension.MergeMextWithExtensions(mext, r.r.Content[mediaType].Extensions)
		if err != nil {
			return nil, err
		}
		media := r.r.Content[mediaType]
		// if mediaExample is the highest priority
		if mext.PayloadGenerationModes[0] == extension.ModeMediaExample {
			// if there's a selector script, it means the user wants to decide which, across multiple examples,
			// should be served
			if mext.MediaExampleSelectorScript != "" {
				return r.runMediaTypeExampleScript(media, mext.MediaExampleSelectorScript)
			} else {
				// otherwise we pick the default example, or any example if the default is missing
				if res := r.runDefaultMediaTypeExample(media); res != nil {
					return res, nil
				}
			}
		}
		if media.Schema == nil {
			return nil, nil
		}
		return jsf.GenerateDataFromSchema(media.Schema.Value, mext, r.evaluator)
	}
	return nil, nil
}

// generateHeaders will generate the response headers using the designated strategy
func (r ResponseDef) generateHeaders(mext extension.Mext) (http.Header, error) {
	headers := http.Header{}
	for k, h := range r.r.Headers {
		if util.RequiredOrRandom(h.Value.Required) {
			val, err := jsf.GenerateDataFromSchema(h.Value.Schema.Value, mext, r.evaluator)
			if err != nil {
				return headers, err
			}
			headers.Set(k, fmt.Sprintf("%v", val))
		}
	}
	return headers, nil
}

// GenerateResponse generates all the necessary pieces of data that comprise a proper response, like payload and
// headers
func (r ResponseDef) GenerateResponse(mext extension.Mext) (*util.Response, error) {
	res := util.NewResponse()
	var err error
	if res.Headers, err = r.generateHeaders(mext); err != nil {
		return res, err
	}
	res.ContentType = r.determineJsonMediaType()
	res.Headers.Set("Content-Type", res.ContentType)
	res.Status = r.status
	if res.Payload, err = r.generateResponsePayload(mext); err != nil {
		return res, err
	}
	return res, nil
}

// runMediaTypeExampleScript executes the media-type example selector script, and retrieves the selected example
func (r ResponseDef) runMediaTypeExampleScript(media *openapi3.MediaType, script string) (any, error) {
	sel, err := r.evaluator.RunScript(script)
	if err != nil {
		return nil, err
	}
	switch t := sel.(type) {
	case string:
		ex := media.Examples[t]
		if ex != nil && ex.Value != nil {
			return ex.Value.Value, nil
		} else {
			return nil, exceptions.Wrap("media_example", errors.New("example "+sel.(string)+" not found"))
		}
	default:
		return nil, exceptions.Wrap("media_example", errors.New("type conversion error"))
	}
}

// runDefaultMediaTypeExample retrieves any media-type example and returns it
func (r ResponseDef) runDefaultMediaTypeExample(media *openapi3.MediaType) any {
	if media.Example != nil {
		return media.Example
	}
	if media.Examples != nil {
		for _, v := range media.Examples {
			return v.Value.Value
		}
	}
	return nil
}
