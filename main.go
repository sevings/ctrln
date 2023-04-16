package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zpatrick/go-config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	toml := config.NewTOMLFile("./config.toml")
	loader := config.NewOnceLoader(toml)
	conf := config.NewConfig([]config.Provider{loader})
	if err := conf.Load(); err != nil {
		log.Fatal(err)
	}

	mqttAddr, err := conf.String("mqtt")
	if err != nil {
		log.Fatal(err)
	}

	mqttToWs := newProxy(mqttAddr, "sevings")
	mqttToWs.SetFormatter(SensorFormatter{})
	err = mqttToWs.Connect()
	if err != nil {
		panic(err)
	}

	err = mqttToWs.Subscribe("sensors/#")
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.Static("/assets/", "./web/assets/")
	router.StaticFile("/", "./web/sensors.html")
	router.GET("/ws", mqttToWs.WsHandler())

	addr, err := conf.String("http")
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Println("Serving web at " + addr)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Exit server")
}
