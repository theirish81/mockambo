package evaluator

import (
	"github.com/cbroglie/mustache"
	"github.com/dop251/goja"
	"mockambo/util"
	"os"
)

type Evaluator struct {
	vm  *goja.Runtime
	ctx map[string]any
}

func NewEvaluator() Evaluator {
	ev := Evaluator{vm: goja.New(), ctx: make(map[string]any)}
	_ = ev.vm.Set("error", "")
	_ = ev.vm.Set("load", ev.Load)
	return ev
}

func (e *Evaluator) Set(key string, val any) {
	e.ctx[key] = val
	_ = e.vm.Set(key, val)
}

func (e *Evaluator) WithRequest(req *util.Request) {
	e.Set("url", req.Request().URL.String())
	e.Set("query", req.Request().URL.Query())
	e.Set("path", req.Request().URL.Path)
	e.Set("method", req.Method)
}

func (e *Evaluator) RunString(script string) (any, error) {
	v, err := e.vm.RunString(script)
	if err != nil {
		return nil, err
	}
	val := v.Export()
	if val, ok := val.(string); ok {
		return e.Template(val)
	}
	return val, err
}

func (e *Evaluator) Load(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (e *Evaluator) Template(templ string) (string, error) {
	tmpl, err := mustache.ParseString(templ)
	if err != nil {
		return "", err
	}
	return tmpl.Render(e.ctx)
}
