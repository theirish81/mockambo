package evaluator

import (
	"github.com/dop251/goja"
	"mockambo/util"
	"os"
)

type Evaluator struct {
	vm *goja.Runtime
}

func NewEvaluator() Evaluator {
	ev := Evaluator{vm: goja.New()}
	_ = ev.vm.Set("error", "")
	_ = ev.vm.Set("load", ev.Load)
	return ev
}

func (e *Evaluator) Set(key string, val any) {
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
	return v.Export(), err
}

func (e *Evaluator) Load(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
