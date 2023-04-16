package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Sensor struct {
	Group  string `json:"group"`
	Sensor string `json:"sensor"`
	Status string `json:"status"`
}

type SensorFormatter struct {
}

func (sf SensorFormatter) FormatMessage(topic string, payload []byte) []byte {
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

	return upd
}
