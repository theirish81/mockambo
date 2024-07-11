package oapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dop251/goja"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"log"
	"mockambo/db"
	"mockambo/extension"
	"mockambo/proxy"
	"mockambo/util"
	"time"
)

type RouteDef struct {
	doc                    Doc
	route                  *routers.Route
	pathItems              map[string]string
	requestValidationInput *openapi3filter.RequestValidationInput
	validationError        error
	mext                   extension.Mext
	vm                     *goja.Runtime
}

func NewRoute(doc Doc, route *routers.Route, pathItems map[string]string) (RouteDef, error) {
	mext, err := extension.MergeDefaultMextWithExtensions(doc.defaultMext, route.Operation.Extensions)
	vm := goja.New()
	_ = vm.Set("pathItems", pathItems)
	_ = vm.Set("error", "")
	return RouteDef{doc: doc, route: route, pathItems: pathItems, mext: mext, vm: vm}, err
}

func (r *RouteDef) OperationID() string {
	return r.route.Operation.OperationID
}

func (r *RouteDef) Process(ctx context.Context, req *util.Request) (*util.Response, error) {
	util.UpdateVmWithRequest(req, r.vm)
	if r.mext.ValidateRequest {
		r.setValidationError(r.validateRequest(ctx, req))

	}
	res := &util.Response{}
	var err error
	if r.mext.Playback {
		key, err := r.vm.RunString(r.mext.RecordingKey)
		if err != nil {
			return res, err
		}
		if data, err := db.Get(key.String(), r.mext.RecordingPath); err == nil {
			log.Println("serving recorded content for key:", key.String())
			err = json.Unmarshal(data, res)
			if err != nil {
				return res, err
			}
			res.Payload, err = base64.StdEncoding.DecodeString(res.Payload.(string))
			sleepTime, err := util.ComputeLatency(r.mext, req)
			if err != nil {
				return res, err
			}
			time.Sleep(sleepTime)
			return res, err
		}
	}
	if r.mext.Proxy {
		if res, err = proxy.Proxy(req.Request(), r.doc.Servers(), r.doc.Servers()[r.mext.ProxyServerIndex]); err != nil {
			return res, err
		}
	} else {
		def, err := r.selectResponse()
		if def == nil || err != nil {
			return nil, err
		}
		res, err = def.GenerateResponseBundle(r.mext)
		if err != nil {
			return res, err
		}
	}
	if r.mext.ValidateResponse {
		if err := r.validateResponse(ctx, res); err != nil {
			return res, err
		}
	}
	_ = r.vm.Set("status", res.Status)
	if r.mext.Record {
		data, err := json.Marshal(res)
		if err != nil {
			return res, err
		}
		key, err := r.vm.RunString(r.mext.RecordingKey)
		if err != nil {
			return res, err
		}
		log.Println("recording content with key:", key.String())
		if err := db.Upsert(key.String(), data, r.mext.RecordingPath); err != nil {
			return res, err
		}
	}
	sleepTime, err := util.ComputeLatency(r.mext, req)
	if err != nil {
		return res, err
	}
	time.Sleep(sleepTime)
	return res, nil
}

func (r *RouteDef) validateRequest(ctx context.Context, req *util.Request) error {
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    req.Request(),
		PathParams: r.pathItems,
		Route:      r.route,
		Options: &openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				return nil
			},
		},
	}
	r.requestValidationInput = requestValidationInput
	return openapi3filter.ValidateRequest(ctx, requestValidationInput)
}

func (r *RouteDef) setValidationError(err error) {
	r.validationError = err
	if r.validationError != nil {
		_ = r.vm.Set("error", "validation_error")
	}
}

func (r *RouteDef) validateResponse(ctx context.Context, bundle *util.Response) error {
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: r.requestValidationInput,
		Status:                 bundle.Status,
		Header:                 bundle.Headers,
	}
	var data []byte
	if dx, ok := bundle.Payload.([]byte); ok {
		data = dx
	} else {
		data, _ = json.Marshal(bundle.Payload)
	}
	responseValidationInput.SetBodyBytes(data)
	return openapi3filter.ValidateResponse(ctx, responseValidationInput)
}

func (r *RouteDef) selectResponse() (*ResponseDef, error) {
	status := 200
	// if ResponseSelector is nil, it implies that a validation error is not managed by a script.
	// Therefore, if a validation error is present, we need to stop the process and let the user know that
	// the request cannot be processed because of that.
	// If, on the contrary, the response selector is a string, therefore a script, then it means that the validation
	// error MAY be handled
	if r.mext.ResponseSelector != nil {
		val, err := r.vm.RunString(*r.mext.ResponseSelector)
		if err != nil {
			return nil, err
		}
		status = int(val.ToInteger())
	}
	def, err := NewResponseDef(r.route.Operation.Responses.Value(fmt.Sprintf("%d", status)).Value, status, r.mext, r.vm)
	return &def, err
}
