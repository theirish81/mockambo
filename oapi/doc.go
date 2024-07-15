package oapi

import (
	"github.com/fsnotify/fsnotify"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"log"
	"mockambo/extension"
	"mockambo/util"
)

type Doc struct {
	t           *openapi3.T
	router      routers.Router
	defaultMext extension.Mext
	watcher     *fsnotify.Watcher
	docPath     string
}

func NewDoc(docPath string) (*Doc, error) {
	doc := Doc{docPath: docPath}

	err := doc.Load()

	return &doc, err
}

func (d *Doc) Load() error {
	log.Println("loading OpenAPI File: ", d.docPath)
	t, err := openapi3.NewLoader().LoadFromFile(d.docPath)
	if err != nil {
		return err
	}
	r, err := gorillamux.NewRouter(t)
	if err != nil {
		return err
	}
	m, err := extension.NewDefaultMextFromExtensions(t.Extensions)
	if err != nil {
		return err
	}
	d.t = t
	d.router = r
	d.defaultMext = m
	return nil
}

func (d *Doc) FindRoute(request *util.Request) (RouteDef, error) {
	r, p, err := d.router.FindRoute(request.Request())
	request.PathItems = p
	if err != nil {
		return RouteDef{}, err
	}
	return NewRouteDef(d, r, p)
}

func (d *Doc) Servers() []string {
	servers := make([]string, 0)
	for _, s := range d.t.Servers {
		servers = append(servers, s.URL)
	}
	return servers
}

func (d *Doc) Watch() error {
	var err error
	d.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case event, ok := <-d.watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if err := d.Load(); err != nil {
						log.Println("could not parse OpenAPI document, running on previous version:", err)
					}
				}
			case err, ok := <-d.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	return d.watcher.Add(d.docPath)
}
