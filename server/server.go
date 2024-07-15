package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"mockambo/oapi"
	"mockambo/util"
	"net/http"
)

type Server struct {
	e    *echo.Echo
	port int
	doc  *oapi.Doc
}

func NewServer(port int, doc *oapi.Doc) Server {
	log.Println("initializing server on port:", port)
	server := Server{e: echo.New(), port: port, doc: doc}
	server.e.HideBanner = true
	server.e.HidePort = true
	server.e.Any("/**", server.handler)
	server.e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Response().Header().Set("X-Mockambo", "true")
		_ = c.JSON(http.StatusInternalServerError, map[string]any{"message": err.Error()})
	}
	return server
}

func (s Server) handler(ctx echo.Context) error {
	util.EnrichURL(ctx)
	req := util.NewRequest(ctx.Request())
	route, err := s.doc.FindRoute(req)
	if err != nil {
		return err
	}
	res, err := route.Process(ctx.Request().Context(), req)
	if err != nil {
		return err
	}
	ctx.Response().Header().Set(util.HeaderContentType, res.ContentType)
	for k, _ := range res.Headers {
		ctx.Response().Header().Set(k, res.Headers.Get(k))
	}
	return util.WriteJSON(ctx, res)
}

func (s Server) Run() error {
	log.Println("starting server")
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}
