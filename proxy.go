package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
)

type Sensor struct {
	Group  string `json:"group"`
	Sensor string `json:"sensor"`
	Status string `json:"status"`
}

type proxy struct {
	updChans map[string]chan []byte
	guard    sync.Mutex
	mqttCli  mqtt.Client
}

func newProxy(address, id string) *proxy {
	p := &proxy{
		updChans: make(map[string]chan []byte),
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(address)
	opts.SetClientID(id)
	opts.SetDefaultPublishHandler(p.mqttHandler())
	opts.SetCleanSession(true)

	p.mqttCli = mqtt.NewClient(opts)

	return p
}

func (p *proxy) addChan(id string) chan []byte {
	updates := make(chan []byte, 10)

	p.guard.Lock()
	defer p.guard.Unlock()

	p.updChans[id] = updates

	return updates
}

func (p *proxy) rmChan(id string) {
	p.guard.Lock()
	defer p.guard.Unlock()

	delete(p.updChans, id)
}

func (p *proxy) Connect() error {
	token := p.mqttCli.Connect()
	token.Wait()
	return token.Error()
}

func (p *proxy) Subscribe(topic string) error {
	token := p.mqttCli.Subscribe(topic, 0, nil)
	token.Wait()
	return token.Error()
}

func (p *proxy) mqttHandler() func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := strings.Split(msg.Topic(), "/")
		if len(topic) != 3 {
			return
		}

		sensor := Sensor{
			Group:  topic[1],
			Sensor: topic[2],
			Status: string(msg.Payload()),
		}
		upd, err := json.Marshal(sensor)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(upd))

		p.guard.Lock()
		defer p.guard.Unlock()

		for _, updates := range p.updChans {
			updates <- upd
		}
	}
}

func (p *proxy) WsHandler() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var upgrader = websocket.Upgrader{}
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		defer func() {
			err := conn.Close()
			if err != nil {
				log.Println(err)
			}
		}()

		addr := ctx.Request.RemoteAddr

		updates := p.addChan(addr)
		defer p.rmChan(addr)

		for upd := range updates {
			err = conn.WriteMessage(websocket.TextMessage, upd)
			if err != nil {
				log.Println(err)
				break
			}
		}
	}
}
