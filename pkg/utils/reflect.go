package utils

import "reflect"

// FieldValueExists uses reflection to check if any fields with the given name have the given value
// The is heavily adapted from https://stackoverflow.com/a/38407429
func FieldValueExists[T comparable](v interface{}, name string, expectedValue T) bool {
	queue := []reflect.Value{reflect.ValueOf(v)}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		// ignore if this is not a struct
		if v.Kind() != reflect.Struct {
			continue
		}
		// iterate through fields looking for match on name
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).Name == name {
				// found it!

				if val, ok := v.Field(i).Interface().(T); ok {
					return val == expectedValue
				}
			}
			// push field to queue
			queue = append(queue, v.Field(i))
		}
	}
	return false
}
