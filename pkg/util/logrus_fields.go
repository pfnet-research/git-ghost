package util

import (
	"reflect"

	log "github.com/Sirupsen/logrus"
)

func ToFields(structObj interface{}) (fields log.Fields) {
	fields = make(log.Fields)

	v := reflect.ValueOf(structObj)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		value := v.FieldByName(name).Interface()
		fields[name] = value
	}

	return
}

func ToFieldsMulti(structObjs ...interface{}) (fields log.Fields) {
	fields = make(log.Fields)
	for structObj := range structObjs {
		fs := ToFields(structObj)
		for k, v := range fs {
			fields[k] = v
		}
	}
	return
}
