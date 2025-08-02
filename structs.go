package vconfig

import (
	"reflect"
)

// isStruct checks if the given interface is a struct type
func isStruct(v any) bool {
	if v == nil {
		return false
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct
}

// structInfo provides information about a struct
type structInfo struct {
	value reflect.Value
	typ   reflect.Type
}

// newStructInfo creates a new structInfo from the given interface
func newStructInfo(v any) *structInfo {
	val := reflect.ValueOf(v)
	typ := val.Type()

	// Dereference pointer if needed
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}

	return &structInfo{
		value: val,
		typ:   typ,
	}
}

// Name returns the name of the struct type
func (s *structInfo) Name() string {
	return s.typ.Name()
}

// FieldOk returns field information if the field exists
func (s *structInfo) FieldOk(name string) (*fieldInfo, bool) {
	field, ok := s.typ.FieldByName(name)
	if !ok {
		return nil, false
	}

	return &fieldInfo{
		field: field,
		value: s.value.FieldByName(name),
	}, true
}

// fieldInfo provides information about a struct field
type fieldInfo struct {
	field reflect.StructField
	value reflect.Value
}

// Kind returns the kind of the field
func (f *fieldInfo) Kind() reflect.Kind {
	return f.field.Type.Kind()
}
