package oapi

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"mockambo/extension"
	"mockambo/util"
)

type Doc struct {
	t           *openapi3.T
	router      routers.Router
	defaultMext extension.Mext
}

func NewDoc(data []byte) (Doc, error) {
	doc := Doc{}
	var err error
	if doc.t, err = openapi3.NewLoader().LoadFromData(data); err != nil {
		return doc, err
	}
	if doc.router, err = gorillamux.NewRouter(doc.t); err != nil {
		return doc, err
	}
	doc.defaultMext, err = extension.NewDefaultMextFromExtensions(doc.t.Extensions)
	return doc, err
}

func (d Doc) FindRoute(request *util.Request) (RouteDef, error) {
	r, p, err := d.router.FindRoute(request.Request())
	request.PathItems = p
	if err != nil {
		return RouteDef{}, err
	}
	return NewRouteDef(d, r, p)
}

func (d Doc) Servers() []string {
	servers := make([]string, 0)
	for _, s := range d.t.Servers {
		servers = append(servers, s.URL)
	}
	return servers
}
