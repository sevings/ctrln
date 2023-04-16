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
}

func newServer(address string) *server {
	router := gin.New()
	router.Use(gin.Recovery())

	router.Static("/assets/", "./web/assets/")
	router.StaticFile("/", "./web/sensors.html")

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &server{
		router: router,
		server: srv,
	}
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
