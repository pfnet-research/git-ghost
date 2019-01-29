package util

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
)

func ToFields(structObj interface{}) (fields log.Fields) {
	fields = make(log.Fields)

	ptyp := reflect.TypeOf(structObj)  // a reflect.Type
	pval := reflect.ValueOf(structObj) // a reflect.Value

	var typ reflect.Type
	var val reflect.Value
	if ptyp.Kind() == reflect.Ptr {
		typ = ptyp.Elem()
		val = pval.Elem()
	} else {
		typ = ptyp
		val = pval
	}
	for i := 0; i < typ.NumField(); i++ {
		name := typ.Field(i).Name
		value := val.FieldByName(name).Interface()
		fields[name] = value
	}

	return
}

func MergeFields(fieldss ...log.Fields) log.Fields {
	merged := make(log.Fields)
	for _, fields := range fieldss {
		for k, v := range fields {
			merged[k] = v
		}
	}
	return merged
}
