package testingdb

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/script-development/RT-CV/db/dbHelpers"
	"github.com/script-development/RT-CV/db/dbInterfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type filter struct {
	filters bson.M
	empty   bool
}

func newFilter(filters ...bson.M) *filter {
	res := &filter{
		filters: dbHelpers.MergeFilters(filters...),
	}
	if len(res.filters) == 0 {
		res.empty = true
	}

	return res
}

func (f *filter) matches(e dbInterfaces.Entry) bool {
	if f.empty {
		return true
	}

	eRefl := reflect.ValueOf(e).Elem()
	eFieldsMap := mapStruct(eRefl.Type())

	for key, value := range f.filters {
		if strings.HasPrefix("$", key) {
			panic("FIXME implement custom filter MongoDB filter properties")
		}

		field, fieldFound := eFieldsMap[key]
		if !fieldFound {
			return false
		}
		goField := eRefl.FieldByName(field.GoFieldName)

		if goField.Kind() == reflect.Ptr {
			if goField.IsNil() {
				return false
			}
			goField = goField.Elem()
		}

		switch typedValue := value.(type) {
		case string:
			if goField.Kind() != reflect.String {
				return false
			}
			if goField.String() != typedValue {
				return false
			}
		case bool:
			if goField.Kind() != reflect.Bool {
				return false
			}
			if goField.Bool() != typedValue {
				return false
			}
		case int:
			switch goField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if goField.Int() != int64(typedValue) {
					return false
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if int64(goField.Uint()) != int64(typedValue) {
					return false
				}
			default:
				return false
			}
		case primitive.ObjectID:
			goFieldValue, ok := goField.Interface().(primitive.ObjectID)
			if !ok {
				return false
			}
			if goFieldValue != typedValue {
				return false
			}
		default:
			panic(fmt.Sprintf("Unimplemented value filter type: %T, key: %v, value: %#v", value, key, value))
		}
	}

	return true
}

type structField struct {
	GoFieldName string
	DbFieldName string
}

func mapStruct(entry reflect.Type) map[string]structField {
	if entry.Kind() != reflect.Struct {
		panic("expected struct but got " + entry.Kind().String())
	}

	res := map[string]structField{}
	for i := 0; i < entry.NumField(); i++ {
		field := entry.Field(i)

		bson := field.Tag.Get("bson")
		if bson == "" {
			bson = field.Tag.Get("json")
		}
		values := strings.Split(bson, ",")
		dbName := values[0]
		if dbName == "" {
			dbName = convertGoToDbName(field.Name)
		}

		res[dbName] = structField{
			GoFieldName: field.Name,
			DbFieldName: dbName,
		}
	}
	return res
}

func convertGoToDbName(fieldname string) string {
	// No need to check if filename length is > 0 beaucase go field name always have a name
	return string(unicode.ToLower(rune(fieldname[0]))) + fieldname[1:]
}
