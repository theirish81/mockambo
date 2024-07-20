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

// Doc is an instrumented OpenAPI Document data structure
type Doc struct {
	t       *openapi3.T
	router  routers.Router
	mext    extension.Mext
	watcher *fsnotify.Watcher
	docPath string
}

// NewDoc will create a new Doc based on the path provided. The path must lead to an OpenAPI spec file
func NewDoc(docPath string) (*Doc, error) {
	doc := Doc{docPath: docPath}

	err := doc.Load()

	return &doc, err
}

// Load will load the OpenAPI specification from the file and initialize all the instrumentation.
// This method can be called multiple times
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
	m, err := extension.NewMextFromExtensions(t.Extensions)
	if err != nil {
		return err
	}
	d.t = t
	d.router = r
	d.mext = m
	return nil
}

// FindRoute will find the route definition in the OpenAPI file, based on the request
func (d *Doc) FindRoute(request *util.Request) (RouteDef, error) {
	r, p, err := d.router.FindRoute(request.Request())
	request.PathItems = p
	if err != nil {
		return RouteDef{}, err
	}
	return NewRouteDef(d, r, p)
}

// Servers will return the list of the servers URLs
func (d *Doc) Servers() []string {
	servers := make([]string, 0)
	for _, s := range d.t.Servers {
		servers = append(servers, s.URL)
	}
	return servers
}

// Watch will start the watching routine on the file described by docPath. If the file changes, then Load
// will be invoked. If Load fails, the Doc data structure will remain unchanged but a log will be printed
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
