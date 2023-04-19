package main

import (
	"encoding/json"
	"sync"
)

type SensorNameStorage struct {
	names map[string]map[string]string
	guard sync.RWMutex
}

func NewSensorNameStorage() *SensorNameStorage {
	return &SensorNameStorage{
		names: make(map[string]map[string]string),
	}
}

func (n *SensorNameStorage) SetName(group, sensor, name string) {
	n.guard.Lock()
	defer n.guard.Unlock()

	groupMap, ok := n.names[group]
	if !ok {
		groupMap = make(map[string]string)
		n.names[group] = groupMap
	}

	groupMap[sensor] = name
}

func (n *SensorNameStorage) GetJSON() []byte {
	data, err := json.Marshal(n.names)
	if err != nil {
		panic(err)
	}

	return data
}
