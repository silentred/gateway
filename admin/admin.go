package admin

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/silentred/glog"
	"github.com/silentred/gateway/config"
	"github.com/silentred/gateway/route"
	"github.com/silentred/gateway/util"
)

var (
	NoError    = util.NewError(0, "ok")
	ParamError = util.NewError(405, "parameter is invalid")
)

// Start listen admin web interface
func Start(cfg *config.Config, table *route.Table, back config.Backend) {
	server := &adminServer{
		table: table,
		back:  back,
	}
	e := echo.New()
	//e.Logger.SetOutput(nil)

	e.Use(middleware.CORS())
	// list services
	e.GET("/table/services", server.listServices)
	// get service
	e.GET("/table/service/:name", server.getService)
	// list routes
	e.GET("/table/routes", server.listRoutes)
	// create/update service
	e.PUT("/table/service", server.putConsulService)
	// del service
	e.DELETE("/table/target/:id", server.delTarget)

	log.Fatal(e.Start(cfg.Admin.Listen))
}

// StartWebUI to serve static webui based on Vue.js
func StartWebUI(cfg *config.Config) {
	e := echo.New()
	//e.Static("/static", "webui/dist/static")
	e.GET("/", func(ctx echo.Context) error {
		b, err := Asset("index.html")
		if err != nil {
			return ctx.String(500, err.Error())
		}
		return ctx.HTML(200, string(b))
	})

	e.GET("/static/*", func(ctx echo.Context) error {
		// strip first char /
		var err error
		var b []byte
		var path = ctx.Request().URL.Path[1:]
		var contentType = "text/html"
		if strings.HasSuffix(path, "css") {
			contentType = "text/css; charset=utf-8"
		}
		if strings.HasSuffix(path, "js") {
			contentType = "application/javascript"
		}

		b, err = Asset(path)
		if err != nil {
			return ctx.String(http.StatusNotFound, err.Error())
		}
		return ctx.Blob(http.StatusOK, contentType, b)
	})

	log.Fatal(e.Start(cfg.WebUI.Listen))
}

type routesDTO struct {
	Route   route.Route    `json:"route"`
	Service *route.Service `json:"service"`
}

type adminServer struct {
	table *route.Table
	back  config.Backend
}

func (server *adminServer) listServices(ctx echo.Context) error {
	svcMap := server.table.Services()
	var svcList = make([]*route.Service, 0, len(svcMap))

	for _, val := range svcMap {
		svcList = append(svcList, val)
	}

	ctx.JSON(http.StatusOK, svcList)
	return nil
}

func (server *adminServer) getService(ctx echo.Context) error {
	var name = ctx.Param("name")

	if name != "" {
		svcMap := server.table.Services()
		if svc, has := svcMap[name]; has {
			return ctx.JSON(http.StatusOK, svc)
		}
	}

	return ctx.JSON(http.StatusNotFound, util.NewError(404, "service not found"))
}

func (server *adminServer) listRoutes(ctx echo.Context) error {
	var dto = []routesDTO{}
	routes := server.table.Routes()

	for r, s := range routes {
		dto = append(dto, routesDTO{
			Route:   r,
			Service: s,
		})
	}

	return ctx.JSON(http.StatusOK, dto)
}

// put service to Consul;
func (server *adminServer) putConsulService(ctx echo.Context) error {
	var err error
	var t config.TargetInfo
	var h config.HealthCheck
	ctx.Bind(&struct {
		*config.TargetInfo
		*config.HealthCheck
	}{&t, &h})

	glog.Debugf("target:%+v , health: %+v", t, h)

	err = server.back.SetTarget(t, h)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.NewError(500, "error on registering service"))
		glog.Error(err)
		return err
	}

	return ctx.JSON(http.StatusOK, NoError)
}

func (server *adminServer) delTarget(ctx echo.Context) error {
	var err error
	var svcID = ctx.Param("id")
	var t config.TargetInfo
	t.ServiceName = ctx.QueryParam("service_name")
	t.TargetHost = ctx.QueryParam("target_host")
	glog.Debugf("del target: Target:%+v ID:%s", t, svcID)

	err = server.back.DelTarget(t, svcID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.NewError(500, "error on deregistering service"))
		glog.Error(err)
		return err
	}

	return ctx.JSON(http.StatusOK, NoError)
}
