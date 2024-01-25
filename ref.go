package gorms

import (
	"reflect"
	"strings"

	"github.com/Drelf2018/TypeGo/Reflect"
)

func contains(s string, substr string) bool {
	return strings.Contains(strings.ToLower(s), substr)
}

func primaryKey(fields []reflect.StructField) string {
	for _, field := range fields {
		if v, ok := field.Tag.Lookup("gorm"); ok {
			if contains(v, "primarykey") {
				return field.Name
			}
		}
	}
	return ""
}

func belongsTo(elem reflect.Type, field reflect.StructField) (exists bool) {
	var name string
	if v, ok := field.Tag.Lookup("gorm"); ok {
		for _, item := range strings.Split(v, ";") {
			if contains(item, "foreignkey") {
				_, name, _ = strings.Cut(item, ":")
				break
			}
		}
	} else {
		key := primaryKey(Reflect.FieldOf(field.Type))
		if key != "" {
			name = field.Name + key
		}
	}
	_, exists = elem.FieldByName(name)
	return
}

func hasOne(key string, field reflect.Type) bool {
	for _, f := range Reflect.FieldOf(field) {
		if f.Name == key {
			return true
		}
	}
	return false
}

type Parser int

func (Parser) Parse(ref *Reflect.Map[[]string], elem reflect.Type) (r []string) {
	fields := Reflect.FieldOf(elem)
	key := elem.Name() + primaryKey(fields)
	for _, field := range fields {
		typ := field.Type
		if typ.Kind() == reflect.Slice {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {
			continue
		}
		if !hasOne(key, typ) && !belongsTo(elem, field) {
			continue
		}
		preloads := ref.MustGetType(typ)
		if len(preloads) == 0 {
			r = append(r, field.Name)
			continue
		}
		prefix := field.Name + "."
		for _, preload := range preloads {
			r = append(r, prefix+preload)
		}
	}
	return
}

var Ref = Reflect.NewMap[Parser, []string](114514, Reflect.SLICEPTRALIAS)
