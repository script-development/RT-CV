package testingdb

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type filter struct {
	filters reflect.Value
	empty   bool
}

func newFilter(filters bson.M) *filter {
	res := &filter{
		filters: reflect.ValueOf(filters),
	}
	if len(filters) == 0 {
		res.empty = true
	}
	return res
}

func (f *filter) matches(value interface{}) bool {
	if f.empty {
		return true
	}

	valueReflection := reflect.ValueOf(value)
	if valueReflection.Kind() == reflect.Ptr {
		valueReflection = valueReflection.Elem()
	}

	return filterMatchesValue(f.filters, valueReflection)
}

func filterMatchesValue(filterMap reflect.Value, value reflect.Value) bool {
	valueFieldsMap, valueIsStruct := mapStruct(value.Type())

	iter := filterMap.MapRange()

filtersLoop:
	for iter.Next() {
		// FIXME we assume the key is a string, maybe we should support also other values
		key := iter.Key().String()
		filter := iter.Value()
		if filter.Kind() == reflect.Interface {
			filter = filter.Elem()
		}

		if strings.HasPrefix(key, "$") {
			panic("FIXME implement custom filter MongoDB filter properties")
		}

		if !valueIsStruct {
			return false
		}

		field, fieldFound := valueFieldsMap[key]
		if !fieldFound {
			return false
		}

		tempValueCopy := value
		for _, goPathPart := range field.GoPathToField {
			tempValueCopy = tempValueCopy.FieldByName(goPathPart)
		}
		valueField := tempValueCopy.FieldByName(field.GoFieldName)

		if valueField.Kind() == reflect.Ptr {
			if valueField.IsNil() {
				return false
			}
			valueField = valueField.Elem()
		}

		if !filter.IsValid() {
			// filter is probably a nil interface{}
			// note that isNil panics if the value is a nil interface without a type
			// so we check here for: interface{}(nil)
			// and not: interface{}([]string(nil))
			switch valueField.Kind() {
			case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
				if valueField.IsNil() {
					continue filtersLoop
				}
			}
			return false
		}

		switch filter.Kind() {
		case reflect.String:
			if valueField.Kind() != reflect.String || valueField.String() != filter.String() {
				return false
			}
		case reflect.Bool:
			if valueField.Kind() != reflect.Bool || valueField.Bool() != filter.Bool() {
				return false
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if !compareInt64ToReflect(filter.Int(), valueField) {
				return false
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if !compareUint64ToReflect(filter.Uint(), valueField) {
				return false
			}
		case reflect.Map:
			if filter.Type().Key().Kind() != reflect.String {
				panic("TODO support filter type map with non string key")
			}
			filterMatchesValue(filter, valueField)
		default:
			filterValue := filter.Interface()
			if filterObjectID, ok := filterValue.(primitive.ObjectID); ok {
				goFieldValue, ok := valueField.Interface().(primitive.ObjectID)
				if !ok {
					return false
				}
				if goFieldValue != filterObjectID {
					return false
				}
			} else {
				panic(fmt.Sprintf(
					"Unimplemented value filter type: %T, key: %v, value: %#v, reflectionKind: %s",
					filterValue,
					key,
					filterValue,
					filter.Kind(),
				))
			}
		}
	}

	return true
}

type structField struct {
	// incase of a inline field we need to resolve the field within another struct
	GoPathToField []string

	GoFieldName string
	DbFieldName string
}

func mapStruct(entry reflect.Type) (structEntries map[string]structField, isStruct bool) {
	if entry.Kind() != reflect.Struct {
		return nil, false
	}

	res := map[string]structField{}
	for i := 0; i < entry.NumField(); i++ {
		mapStructField(entry.Field(i), func(field structField) {
			res[field.DbFieldName] = field
		})
	}
	return res, true
}

func mapStructField(field reflect.StructField, add func(structField)) {
	bsonTag := field.Tag.Get("bson")
	if bsonTag == "" {
		bsonTag = field.Tag.Get("json")
	}

	values := strings.Split(bsonTag, ",")
	dbName := values[0]
	if dbName == "" {
		dbName = convertGoToDbName(field.Name)
	}

	isInlineField := false
	if len(values) > 1 {
		for _, entry := range values[1:] {
			if entry == "inline" && field.Type.Kind() == reflect.Struct {
				isInlineField = true
			}
		}
	}

	if isInlineField {
		for i := 0; i < field.Type.NumField(); i++ {
			mapStructField(field.Type.Field(i), func(toAdd structField) {
				toAdd.GoPathToField = append(toAdd.GoPathToField, field.Name)
				add(toAdd)
			})
		}
	} else {
		add(structField{
			GoPathToField: []string{},
			GoFieldName:   field.Name,
			DbFieldName:   dbName,
		})
	}
}

func convertGoToDbName(fieldname string) string {
	// No need to check if filename length is > 0 beaucase go field name always have a name
	return string(unicode.ToLower(rune(fieldname[0]))) + fieldname[1:]
}

func compareInt64ToReflect(value int64, reflectionValue reflect.Value) bool {
	switch reflectionValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflectionValue.Int() == value
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflectionIntValue := int64(reflectionValue.Uint())
		if reflectionIntValue < 0 {
			// The uint64 value of the reflect value was more than the highest int64 value and thus resetted itself and now it's below zero
			return false
		}
		return reflectionIntValue == value
	default:
		return false
	}
}

func compareUint64ToReflect(value uint64, reflectionValue reflect.Value) bool {
	switch reflectionValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue := reflectionValue.Int()
		if intValue < 0 {
			return false
		}
		return uint64(intValue) == value
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflectionValue.Uint() == value
	default:
		return false
	}
}
