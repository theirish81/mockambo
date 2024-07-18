package evaluator

import (
	"github.com/cbroglie/mustache"
	"github.com/dop251/goja"
	"mockambo/exceptions"
	"mockambo/util"
	"os"
)

// Evaluator is the script and template evaluator
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

// Set sets a variable in the evaluator scope
func (e *Evaluator) Set(key string, val any) {
	e.ctx[key] = val
	_ = e.vm.Set(key, val)
}

// WithRequest extracts important values from a util.Request and sets them in the scope of the evaluator
func (e *Evaluator) WithRequest(req *util.Request) {
	e.Set("url", req.Request().URL.String())
	e.Set("query", req.Request().URL.Query())
	e.Set("path", req.Request().URL.Path)
	e.Set("method", req.Method)
}

// RunScript evaluates a JavaScript script
func (e *Evaluator) RunScript(script string) (any, error) {
	v, err := e.vm.RunString(script)
	if err != nil {
		return nil, exceptions.Wrap("evaluate", err)
	}
	val := v.Export()
	if val, ok := val.(string); ok {
		t, err := e.Template(val)
		if err != nil {
			err = exceptions.Wrap("template", err)
		}
		return t, err
	}
	return val, nil
}

// Load loads a text file. This function gets injected in the JavaScript scope
func (e *Evaluator) Load(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", exceptions.Wrap("load", err)
	}
	return string(data), nil
}

// Template renders a template against the Evaluator scope
func (e *Evaluator) Template(templ string) (string, error) {
	tmpl, err := mustache.ParseString(templ)
	if err != nil {
		return "", err
	}
	return tmpl.Render(e.ctx)
}
