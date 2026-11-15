package pgsql

import (
	"reflect"
	"strings"
	"sync"
)

type ColumnMapping struct {
	Column string
	Index  int
	PK     bool
}

type TableMapping struct {
	Columns []ColumnMapping
	PKIndex int
}

var mappingCache sync.Map

func MappingFor(t reflect.Type) TableMapping {
	if cached, ok := mappingCache.Load(t); ok {
		return cached.(TableMapping)
	}
	m := buildMapping(t)
	mappingCache.Store(t, m)
	return m
}

func buildMapping(t reflect.Type) TableMapping {
	m := TableMapping{PKIndex: -1}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("store")
		if tag == "" || tag == "-" {
			continue
		}
		parts := strings.Split(tag, ",")
		cm := ColumnMapping{Column: parts[0], Index: i}
		for _, p := range parts[1:] {
			if p == "pk" {
				cm.PK = true
			}
		}
		if cm.PK {
			m.PKIndex = len(m.Columns)
		}
		m.Columns = append(m.Columns, cm)
	}
	return m
}
