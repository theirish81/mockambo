package oapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"log"
	"mockambo/db"
	"mockambo/evaluator"
	"mockambo/exceptions"
	"mockambo/extension"
	"mockambo/jsf"
	"mockambo/proxy"
	"mockambo/util"
	"strconv"
	"time"
)

type RouteDef struct {
	doc                    *Doc
	route                  *routers.Route
	pathItems              map[string]string
	requestValidationInput *openapi3filter.RequestValidationInput
	validationError        error
	mext                   extension.Mext
	evaluator              evaluator.Evaluator
}

func NewRouteDef(doc *Doc, route *routers.Route, pathItems map[string]string) (RouteDef, error) {
	mext, err := extension.MergeDefaultMextWithExtensions(doc.defaultMext, route.Operation.Extensions)
	ev := evaluator.NewEvaluator()
	ev.Set("fake", jsf.Fake)
	ev.Set("pathItems", pathItems)
	return RouteDef{doc: doc, route: route, pathItems: pathItems, mext: mext, evaluator: ev}, err
}

func (r *RouteDef) OperationID() string {
	return r.route.Operation.OperationID
}

func (r *RouteDef) Process(ctx context.Context, req *util.Request) (*util.Response, error) {
	r.evaluator.WithRequest(req)
	if r.mext.ValidateRequest {
		r.setValidationError(r.validateRequest(ctx, req))

	}
	res := util.NewResponse()
	var err error
	if r.mext.Playback {
		key, err := r.evaluator.RunString(r.mext.RecordingKey)
		if err != nil {
			return res, exceptions.Wrap("playback", err)
		}
		if data, err := db.Get(key.(string), r.mext.RecordingPath); err == nil {
			log.Println("serving recorded content for key:", key.(string))
			err = json.Unmarshal(data, res)
			if err != nil {
				return res, exceptions.Wrap("playback", err)
			}
			res.Payload, err = base64.StdEncoding.DecodeString(res.Payload.(string))
			sleepTime, err := util.ComputeLatency(r.mext, req)
			if err != nil {
				return res, exceptions.Wrap("playback", err)
			}
			time.Sleep(sleepTime)
			return res, nil
		}
	}
	if r.mext.Proxy {
		if res, err = proxy.Proxy(req.Request(), r.doc.Servers(), r.doc.Servers()[r.mext.ProxyServerIndex]); err != nil {
			return res, exceptions.Wrap("proxy", err)
		}
	} else {
		def, err := r.selectResponse()
		if def == nil || err != nil {
			return nil, exceptions.Wrap("select_response", err)
		}
		res, err = def.GenerateResponseBundle(r.mext)
		if err != nil {
			return res, exceptions.Wrap("generate_response", err)
		}
	}
	if r.mext.ValidateResponse {
		if err := r.validateResponse(ctx, res); err != nil {
			return res, exceptions.Wrap("validate_response", err)
		}
	}
	r.evaluator.Set("status", res.Status)
	if r.mext.Record {
		data, err := json.Marshal(res)
		if err != nil {
			return res, exceptions.Wrap("record", err)
		}
		key, err := r.evaluator.RunString(r.mext.RecordingKey)
		if err != nil {
			return res, exceptions.Wrap("record", err)
		}
		log.Println("recording content with key:", key.(string))
		if err := db.Upsert(key.(string), data, r.mext.RecordingPath); err != nil {
			return res, exceptions.Wrap("record", err)
		}
	}
	sleepTime, err := util.ComputeLatency(r.mext, req)
	if err != nil {
		return res, exceptions.Wrap("latency", err)
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
		r.evaluator.Set("error", "validation_error")
	}
}

func (r *RouteDef) validateResponse(ctx context.Context, bundle *util.Response) error {
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: r.requestValidationInput,
		Status:                 bundle.Status,
		Header:                 bundle.Headers,
	}
	var data []byte
	switch t := bundle.Payload.(type) {
	case []byte:
		data = t
	case string:
		data = []byte(t)
	default:
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
		val, err := r.evaluator.RunString(*r.mext.ResponseSelector)
		if err != nil {
			return nil, err
		}
		status = int(val.(int64))
	}
	selector := fmt.Sprintf("%d", status)
	if res := r.route.Operation.Responses.Value(selector); res != nil {
		def, err := NewResponseDef(res.Value, status, r.mext, r.evaluator)
		return &def, err
	} else if res := r.route.Operation.Responses.Value("default"); res != nil {
		def, err := NewResponseDef(res.Value, status, r.mext, r.evaluator)
		return &def, err
	} else {
		for k, _ := range r.route.Operation.Responses.Map() {
			val, _ := strconv.Atoi(k)
			def, err := NewResponseDef(r.route.Operation.Responses.Value(k).Value, val, r.mext, r.evaluator)
			return &def, err
		}

	}

	return nil, exceptions.Wrap("route_find", errors.New("route not found"))
}
