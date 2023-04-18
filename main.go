package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zpatrick/go-config"
	"log"
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
	mqttToWs.SetFormatter(NewSensorFormatter())
	err = mqttToWs.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = mqttToWs.Subscribe("sensors/#")
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)

	httpAddr, err := conf.String("http")
	if err != nil {
		log.Fatal(err)
	}

	srv := newServer(httpAddr)
	srv.SetWsHandler(mqttToWs.WsHandler())

	//ioAddr, err := conf.String("socket_io")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//ioOpts := &socketio_client.Options{}
	//io, err := socketio_client.NewClient(ioAddr, ioOpts)
	//if err != nil {
	//	log.Fatal("io: ", err)
	//}

	//srv.SetEmitter(io)

	log.Println("Serving web at " + httpAddr)
	go srv.Listen()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Exit server")
}
