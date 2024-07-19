package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"mockambo/exceptions"
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
		log.Println(fmt.Sprintf("%s -> %d, (%s)", util.NewRequest(c.Request()).String(), http.StatusInternalServerError, err.Error()))
		c.Response().Header().Set("X-Mockambo", "true")
		switch t := err.(type) {
		case *exceptions.MockamboError:
			_ = c.JSON(http.StatusInternalServerError, t)
		default:
			_ = c.JSON(http.StatusInternalServerError, map[string]any{"source": "mockambo", "message": err.Error()})
		}
	}
	return server
}

func (s Server) handler(ctx echo.Context) error {
	util.EnrichURL(ctx)
	req := util.NewRequest(ctx.Request())
	route, err := s.doc.FindRoute(req)
	if err != nil {
		return exceptions.Wrap("find_route", err)
	}
	res, err := route.Process(ctx.Request().Context(), req)
	if err != nil {
		return err
	}
	ctx.Response().Header().Set("Content-Type", res.ContentType)
	for k, _ := range res.Headers {
		ctx.Response().Header().Set(k, res.Headers.Get(k))
	}
	log.Println(fmt.Sprintf("%s -> %d", req.String(), res.Status))
	return util.WriteJSON(ctx, res)
}

func (s Server) Run() error {
	log.Println("starting server")
	return s.e.Start(fmt.Sprintf(":%d", s.port))
}
