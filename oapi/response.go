package oapi

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/getkin/kin-openapi/openapi3"
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
	vm          *goja.Runtime
}

func NewResponseDef(r *openapi3.Response, status int, mext extension.Mext, vm *goja.Runtime) (ResponseDef, error) {
	defaultMext, err := extension.MergeDefaultMextWithExtensions(mext, r.Extensions)
	_ = vm.Set("status", status)
	return ResponseDef{r: r, status: status, defaultMext: defaultMext, vm: vm}, err
}

func (r ResponseDef) determineMediaType() string {
	for k, _ := range r.r.Content {
		if strings.Contains(k, "json") {
			return k
		}
	}
	return ""
}

func (r ResponseDef) generateResponsePayload(mext extension.Mext) (any, error) {
	if mediaType := r.determineMediaType(); mediaType != "" {
		return jsf.GenerateDataFromSchema(r.r.Content[mediaType].Schema.Value, mext, r.vm)
	}
	return nil, nil
}

func (r ResponseDef) generateHeaders(mext extension.Mext) (http.Header, error) {
	headers := http.Header{}
	for k, h := range r.r.Headers {
		if util.RequiredOrRandom(h.Value.Required) {
			val, err := jsf.GenerateDataFromSchema(h.Value.Schema.Value, mext, r.vm)
			if err != nil {
				return headers, err
			}
			headers.Set(k, fmt.Sprintf("%v", val))
		}
	}
	return headers, nil
}

func (r ResponseDef) GenerateResponseBundle(mext extension.Mext) (*util.Response, error) {
	res := util.Response{}
	var err error
	if res.Headers, err = r.generateHeaders(mext); err != nil {
		return &res, err
	}
	res.ContentType = r.determineMediaType()
	res.Headers.Set(util.HeaderContentType, res.ContentType)
	res.Status = r.status
	if res.Payload, err = r.generateResponsePayload(mext); err != nil {
		return &res, err
	}
	return &res, nil
}
