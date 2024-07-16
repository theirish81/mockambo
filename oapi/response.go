package oapi

import (
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"mockambo/evaluator"
	"mockambo/extension"
	"mockambo/jsf"
	"mockambo/util"
	"net/http"
	"strings"
)

type ResponseDef struct {
	status      int
	r           *openapi3.Response
	err         error
	defaultMext extension.Mext
	evaluator   evaluator.Evaluator
}

func NewResponseDef(r *openapi3.Response, status int, mext extension.Mext, ev evaluator.Evaluator) (ResponseDef, error) {
	defaultMext, err := extension.MergeDefaultMextWithExtensions(mext, r.Extensions)
	ev.Set("status", status)
	return ResponseDef{r: r, status: status, defaultMext: defaultMext, evaluator: ev}, err
}

func (r ResponseDef) determineMediaType() string {
	for k := range r.r.Content {
		if strings.Contains(k, "json") {
			return k
		}
	}
	return ""
}

func (r ResponseDef) generateResponsePayload(mext extension.Mext) (any, error) {
	if mediaType := r.determineMediaType(); mediaType != "" {
		mext, err := extension.MergeDefaultMextWithExtensions(mext, r.r.Content[mediaType].Extensions)
		if err != nil {
			return nil, err
		}
		if r.r.Content[mediaType].Schema == nil {
			return nil, nil
		}
		return jsf.GenerateDataFromSchema(r.r.Content[mediaType].Schema.Value, mext, r.evaluator)
	}
	return nil, nil
}

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

func (r ResponseDef) GenerateResponseBundle(mext extension.Mext) (*util.Response, error) {
	res := util.NewResponse()
	var err error
	if res.Headers, err = r.generateHeaders(mext); err != nil {
		return res, err
	}
	res.ContentType = r.determineMediaType()
	res.Headers.Set(util.HeaderContentType, res.ContentType)
	res.Status = r.status
	if res.Payload, err = r.generateResponsePayload(mext); err != nil {
		return res, err
	}
	return res, nil
}
