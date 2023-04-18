package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Sensor struct {
	Group    string `json:"group"`
	Sensor   string `json:"sensor"`
	Status   string `json:"status"`
	Critical bool   `json:"critical,omitempty"`
}

type SensorFormatter struct {
	sensors  map[string][]byte
	guard    sync.Mutex
	statusRe *regexp.Regexp
}

func NewSensorFormatter() *SensorFormatter {
	return &SensorFormatter{
		sensors:  make(map[string][]byte),
		statusRe: regexp.MustCompile(`^\W*([\d.]*)\W*(\w*)\W*$`),
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
	}

	value, text, status := sf.convertValue(string(payload))
	sensor.Critical = sf.checkCritical(value, text)
	sensor.Status = status

	upd, err := json.Marshal(sensor)
	if err != nil {
		log.Println(err)
		return nil
	}

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

func (sf *SensorFormatter) convertValue(status string) (float64, string, string) {
	match := sf.statusRe.FindStringSubmatch(status)
	if len(match) < 2 {
		return 0, "", status
	}

	text := match[2]
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0, text, status
	}

	if text == "F" {
		value = math.Round((value - 32) / 1.8)
		text = "C"
	}

	status = fmt.Sprintf("%g %s", value, text)

	return value, text, status
}

func (sf *SensorFormatter) checkCritical(value float64, text string) bool {
	switch text {
	case "Lux":
		return value > 800
	case "C":
		return value > 50
	case "V":
		return value < 170
	case "open":
		return true
	case "closed":
		return false
	case "off":
		return false
	case "on":
		return false
	}

	return true
}
