package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

type server struct {
	router *gin.Engine
	server *http.Server
	emit   Emitter
	names  map[string]map[string]string
	guard  sync.RWMutex
}

type Emitter interface {
	Emit(message string, args ...interface{}) error
}

type ConsoleEmitter struct {
}

func (ce ConsoleEmitter) Emit(message string, args ...interface{}) error {
	log.Println("emit:", message, args)
	return nil
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
		emit:  ConsoleEmitter{},
		names: make(map[string]map[string]string),
	}

	router.Static("/assets/", "./web/assets/")
	router.StaticFile("/", "./web/sensors.html")
	router.POST("/:group/:sensor/name", srv.postSensorHandler())
	router.GET("/sensors", srv.getSensorsHandler())
	router.POST("/:group/on", srv.switchGroupHandler("on"))
	router.POST("/:group/off", srv.switchGroupHandler("off"))

	return srv
}

func (srv *server) SetWsHandler(handler gin.HandlerFunc) {
	srv.router.GET("/ws", handler)
}

func (srv *server) SetEmitter(emit Emitter) {
	srv.emit = emit
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

		srv.guard.Lock()
		defer srv.guard.Unlock()

		groupMap, ok := srv.names[group]
		if !ok {
			groupMap = make(map[string]string)
			srv.names[group] = groupMap
		}

		groupMap[sensor] = name
		ctx.Status(200)
	}
}

func (srv *server) getSensorsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		srv.guard.RLock()
		defer srv.guard.RUnlock()

		ctx.JSON(200, srv.names)
	}
}

type switchGroupMsg struct {
	Group string `json:"group"`
	Value string `json:"value"`
}

func (srv *server) switchGroupHandler(value string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		msg := switchGroupMsg{
			Group: ctx.Param("group"),
			Value: value,
		}

		err := srv.emit.Emit("enable", msg)
		if err != nil {
			log.Println(err)
		}

		ctx.Status(200)
	}
}
