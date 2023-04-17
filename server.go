package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type server struct {
	router *gin.Engine
	server *http.Server
	names  map[string]map[string]string
}

func newServer(address string) *server {
	router := gin.New()
	router.Use(gin.Recovery())

	srv := &server{
		router: router,
		server: &http.Server{
			Addr:    address,
			Handler: router,
		},
		names: make(map[string]map[string]string),
	}

	router.Static("/assets/", "./web/assets/")
	router.StaticFile("/", "./web/sensors.html")
	router.POST("/[group]/[sensor]", srv.postSensorHandler())
	router.GET("/sensors")

	return srv
}

func (srv *server) SetWsHandler(handler gin.HandlerFunc) {
	srv.router.GET("/ws", handler)
}

func (srv *server) Listen() {
	err := srv.server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (srv *server) Shutdown(ctx context.Context) {
	err := srv.server.Shutdown(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (srv *server) postSensorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		group := ctx.Param("group")
		sensor := ctx.Param("sensor")
		name := ctx.PostForm("name")

		groupMap, ok := srv.names[group]
		_, ok = groupMap[sensor]
		if !ok {
			ctx.AbortWithStatus(400)
			return
		}

		groupMap[sensor] = name
		ctx.Status(200)
	}
}

func (srv *server) getSensorsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, srv.names)
	}
}
