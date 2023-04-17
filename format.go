package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
)

type Sensor struct {
	Group  string `json:"group"`
	Sensor string `json:"sensor"`
	Status string `json:"status"`
}

type SensorFormatter struct {
	sensors map[string][]byte
	guard   sync.Mutex
}

func NewSensorFormatter() *SensorFormatter {
	return &SensorFormatter{
		sensors: make(map[string][]byte),
	}
}

func (sf *SensorFormatter) FormatMessage(topic string, payload []byte) []byte {
	names := strings.Split(topic, "/")
	if len(names) != 3 {
		return nil
	}

	sensor := Sensor{
		Group:  names[1],
		Sensor: names[2],
		Status: string(payload),
	}

	upd, err := json.Marshal(sensor)
	if err != nil {
		log.Println(err)
		return nil
	}

	fmt.Println(string(upd))

	sf.guard.Lock()
	defer sf.guard.Unlock()

	sf.sensors[topic] = upd

	return upd
}

func (sf *SensorFormatter) MessagesOnConnect() [][]byte {
	sf.guard.Lock()
	defer sf.guard.Unlock()

	var messages [][]byte
	for _, msg := range sf.sensors {
		messages = append(messages, msg)
	}

	return messages
}
