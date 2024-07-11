package util

import (
	"github.com/dop251/goja"
)

func UpdateVmWithRequest(req *Request, vm *goja.Runtime) {
	_ = vm.Set("url", req.Request().URL.String())
	_ = vm.Set("query", req.Request().URL.Query())
	_ = vm.Set("path", req.Request().URL.Path)
	_ = vm.Set("method", req.Method)
}
