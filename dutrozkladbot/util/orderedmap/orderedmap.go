// Package orderedmap ..
package orderedmap

import (
	"bytes"
	"encoding/json"
	"sort"
)

// OrderedMap ..
type OrderedMap struct {
	Order []int
	Map   map[int]interface{}
}

// New ..
func New() *OrderedMap {
	o := OrderedMap{}

	o.Order = make([]int, 0)
	o.Map = make(map[int]interface{})
	return &o
}

// NewLen ..
func NewLen(len int) *OrderedMap {
	o := OrderedMap{}

	o.Order = make([]int, 0, len)
	o.Map = make(map[int]interface{})
	return &o
}

// Set ..
func (om *OrderedMap) Set(key int, value interface{}) {
	if _, exists := om.Map[key]; !exists {
		om.Order = append(om.Order, key)
	}

	om.Map[key] = &value
}

// Push ..
func (om *OrderedMap) Push(value interface{}) {
	key := len(om.Order) - 1

	om.Order = append(om.Order, key)
	om.Map[key] = value
}

// UnmarshalJSON ..
func (om *OrderedMap) UnmarshalJSON(b []byte) error {
	json.Unmarshal(b, &om.Map)

	index := make(map[int]int)
	for key := range om.Map {
		om.Order = append(om.Order, key)
		esc, _ := json.Marshal(key) //Escape the key
		index[key] = bytes.Index(b, esc)
	}

	sort.Slice(om.Order, func(i, j int) bool { return index[om.Order[i]] < index[om.Order[j]] })
	return nil
}
