package startype

import (
	"fmt"
	"reflect"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type StarValue struct {
	val starlark.Value
}

// Starlark wraps a Starlark value val
// so it can be converted to a Go value.
func Starlark(val starlark.Value) *StarValue {
	return &StarValue{val: val}
}

// Value returns the wrapped Starlark value
func (v *StarValue) Value() starlark.Value {
	return v.val
}

// Go converts Starlark the wrapped value and stores the
// result into a Go value specified by pointer goPtr.
// Example:
//
//    var msg string
//    Star(starlark.String("Hello")).Go(&msg)
//
// This method supports the following type map from Starlark to Go types:
//
//      starlark.Bool   	-- bool
//      starlark.Int    	-- int64 or uint64
//      starlark.Float  	-- float64
//      starlark.String 	-- string
//      *starlark.List  	-- []T
//      starlark.Tuple  	-- []T
//      *starlark.Dict  	-- map[K]T
//      *starlark.Set   	-- []T

func (v *StarValue) Go(goin interface{}) error {
	goval := reflect.ValueOf(goin)
	gotype := goval.Type()
	if gotype.Kind() != reflect.Pointer || goval.IsNil() {
		return fmt.Errorf("Go target must be a poiner or addressable: got %v", gotype)
	}

	return starlarkToGo(v.val, goval.Elem())
}

// starlarkToGo translates starlark.Archive val to the provided Go value goval
// using the following type mapping:
//
//      starlark.Bool   	-- bool
//      starlark.Int    	-- int64 or uint64
//      starlark.Float  	-- float64
//      starlark.String 	-- string
//      *starlark.List  	-- []T
//      starlark.Tuple  	-- []T
//      *starlark.Dict  	-- map[K]T
//      *starlark.Set   	-- []T

func starlarkToGo(srcVal starlark.Value, goval reflect.Value) error {
	if srcVal == nil {
		return nil
	}

	gotype := goval.Type()

	var starval reflect.Value
	srcType := srcVal.Type()

	switch srcType {
	case "bool":
		if gotype.Kind() != reflect.Bool && gotype.Kind() != reflect.Interface && gotype.Kind() != reflect.Pointer {
			return fmt.Errorf("target type (%s): must be bool, *bool, or any", gotype.Kind())
		}
		starval = reflect.ValueOf(bool(srcVal.Truth()))

		if gotype.Kind() == reflect.Pointer {
			goval.Set(reflect.New(gotype.Elem()))
			return starlarkToGo(srcVal, goval.Elem()) // convert using value instead of pointer
		}

		goval.Set(starval)
		return nil

	case "int":
		intVal, ok := srcVal.(starlark.Int)
		if !ok {
			return fmt.Errorf("source value must be starlark.Int: %T", srcVal)
		}

		switch gotype.Kind() {
		case reflect.Pointer:
			goval.Set(reflect.New(gotype.Elem()))
			return starlarkToGo(srcVal, goval.Elem())
		case reflect.Uint8:
			if val, ok := intVal.Int64(); ok {
				starval = reflect.ValueOf(uint8(val))
			}
		case reflect.Int8:
			if val, ok := intVal.Int64(); ok {
				starval = reflect.ValueOf(int8(val))
			}
		case reflect.Uint16:
			if val, ok := intVal.Int64(); ok {
				starval = reflect.ValueOf(uint16(val))
			}
		case reflect.Int16:
			if val, ok := intVal.Int64(); ok {
				starval = reflect.ValueOf(int16(val))
			}
		case reflect.Int:
			if val, ok := intVal.Int64(); ok {
				starval = reflect.ValueOf(int(val))
			}
		case reflect.Int32:
			if val, ok := intVal.Int64(); ok {
				starval = reflect.ValueOf(int32(val))
			}
		case reflect.Uint:
			if val, ok := intVal.Uint64(); ok {
				starval = reflect.ValueOf(uint(val))
			}
		case reflect.Uint32:
			if val, ok := intVal.Uint64(); ok {
				starval = reflect.ValueOf(uint32(val))
			}
		case reflect.Int64, reflect.Uint64, reflect.Interface:
			bigInt := intVal.BigInt()
			switch {
			case bigInt.IsInt64():
				starval = reflect.ValueOf(bigInt.Int64())
			case bigInt.IsUint64():
				starval = reflect.ValueOf(bigInt.Uint64())
			default:
				return fmt.Errorf("unsupported starlark.Int type")
			}
		default:
			return fmt.Errorf("unsupported target type (%v): must be int, int8, int16, int32, uint, uint32, int64, uint64, pointers to them, or any", gotype.Kind())
		}

		goval.Set(starval)
		return nil

	case "float":
		if gotype.Kind() != reflect.Float64 && gotype.Kind() != reflect.Float32 && gotype.Kind() != reflect.Interface && gotype.Kind() != reflect.Pointer {
			return fmt.Errorf("target type (%s): must be float32, float64, *float{32|64}, or any", gotype.Kind())
		}
		floatVal, ok := srcVal.(starlark.Float)
		if !ok {
			return fmt.Errorf("source value must starlark.Float: %T", srcVal)
		}

		switch gotype.Kind() {
		case reflect.Pointer:
			goval.Set(reflect.New(gotype.Elem()))
			return starlarkToGo(srcVal, goval.Elem())
		case reflect.Float32:
			starval = reflect.ValueOf(float32(floatVal))
		case reflect.Float64, reflect.Interface:
			starval = reflect.ValueOf(float64(floatVal))
		default:
			return fmt.Errorf("unsupported float target:: %s", gotype.Kind())
		}

		goval.Set(starval)
		return nil

	case "string":
		if gotype.Kind() != reflect.String && gotype.Kind() != reflect.Interface && gotype.Kind() != reflect.Pointer {
			return fmt.Errorf("Starlark.String to Go: target target (%s): must be string, *string, or any", gotype.Kind())
		}

		strVal, ok := srcVal.(starlark.String)
		if !ok {
			return fmt.Errorf("Starlark.String to Go: failed to assert %T as starlark.String", srcVal)
		}

		starval = reflect.ValueOf(string(strVal))

		if gotype.Kind() == reflect.Pointer {
			goval.Set(reflect.New(gotype.Elem()))
			return starlarkToGo(srcVal, goval.Elem())
		}

		goval.Set(starval)
		return nil

	case "list":
		if gotype.Kind() != reflect.Slice && gotype.Kind() != reflect.Array {
			return fmt.Errorf("target type must be slice or array")
		}
		listVal, ok := srcVal.(*starlark.List)
		if !ok {
			return fmt.Errorf("failed to assert %T as *starlark.List", srcVal)
		}
		goval.Set(reflect.MakeSlice(gotype, listVal.Len(), listVal.Len()))
		for i := 0; i < listVal.Len(); i++ {
			if err := starlarkToGo(listVal.Index(i), goval.Index(i)); err != nil {
				return err
			}
		}
		return nil

	case "tuple":
		if gotype.Kind() != reflect.Slice && gotype.Kind() != reflect.Array {
			return fmt.Errorf("target type must be slice or array")
		}
		tupVal, ok := srcVal.(starlark.Tuple)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Tuple", srcVal)
		}
		goval.Set(reflect.MakeSlice(gotype, tupVal.Len(), tupVal.Len()))
		for i := 0; i < tupVal.Len(); i++ {
			if err := starlarkToGo(tupVal.Index(i), goval.Index(i)); err != nil {
				return err
			}
		}
		return nil

	case "dict":
		// Converting a Dict -> Map requires a bit of work to handle embedded maps,
		// when the outer map value is of type `map[T]any`. The reflect package cannot build
		// a new map without type information for both key and elements. So, extra work must be
		// done to construct an inner map dynamically by assuming type `map[any]any` when the
		// outer map values have type like `map[string]any` for instance.
		dict, ok := srcVal.(*starlark.Dict)
		if !ok {
			return fmt.Errorf("failed to assert %T as *starlark.Dict", srcVal)
		}

		// map target type
		switch gotype.Kind() {
		case reflect.Map:
			goval.Set(reflect.MakeMapWithSize(gotype, dict.Len()))
		case reflect.Interface:
			gotype = reflect.TypeOf(map[any]any{})
			goval.Set(reflect.MakeMapWithSize(gotype, dict.Len()))
		case reflect.Pointer:
			goval.Set(reflect.New(gotype.Elem()))
			return starlarkToGo(dict, goval.Elem())
		default:
			return fmt.Errorf("Starlark.Dict to Go: target type (%s): must be map, *map, any, or pointer", gotype.Name())
		}

		for _, dictKey := range dict.Keys() {
			dictVal, ok, err := dict.Get(dictKey)
			if err != nil {
				return fmt.Errorf("starlark.Dict.Get failed: %s", err)
			}
			if !ok {
				continue
			}

			// convert map key
			keyType := getExactMapType(dictKey, gotype.Key())
			goMapKey := reflect.New(keyType).Elem()
			if err := starlarkToGo(dictKey, goMapKey); err != nil {
				return err
			}

			// convert map element
			var goMapElem reflect.Value
			if dictVal != nil {
				elemType := getExactMapType(dictVal, gotype.Elem())
				goMapElem = reflect.New(elemType).Elem()
				if err := starlarkToGo(dictVal, goMapElem); err != nil {
					return err
				}
			} else {
				goMapElem = reflect.ValueOf(nil)
			}

			// store map value
			goval.SetMapIndex(goMapKey, goMapElem)
		}
		return nil

	case "set":
		if gotype.Kind() != reflect.Slice && gotype.Kind() != reflect.Array {
			return fmt.Errorf("target type must be slice or array")
		}
		setVal, ok := srcVal.(*starlark.Set)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Tuple", srcVal)
		}
		goval.Set(reflect.MakeSlice(gotype, setVal.Len(), setVal.Len()))
		var setItem starlark.Value
		iter := setVal.Iterate()
		i := 0
		for iter.Next(&setItem) {
			if err := starlarkToGo(setItem, goval.Index(i)); err != nil {
				return err
			}
			i++
		}
		return nil

	case "struct":
		if gotype.Kind() != reflect.Struct {
			return fmt.Errorf("target type (%s): must be a struct ", gotype.Kind())
		}

		structVal, ok := srcVal.(*starlarkstruct.Struct)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Struct", srcVal)
		}

		// copy starlark struct attributes to struct fields
		attrs := structVal.AttrNames()
		for _, attr := range attrs {
			attrVal, err := structVal.Attr(attr)
			if err != nil {
				return fmt.Errorf("starlarkstruct.Struct attribute %s: %s", attr, err)
			}

			// determine struct field name from struct tag or starlarkstruct field name attribute
			fieldName, found := findStructFieldByTag(gotype, "name", attr)
			if !found {
				fieldName = strings.Title(attr) //nolint:staticcheck
			}

			// decode struct field
			if field, ok := gotype.FieldByName(fieldName); ok {

				fieldVal := goval.FieldByName(field.Name)

				if fieldVal.Kind() == reflect.Pointer {
					fieldVal.Set(reflect.New(field.Type.Elem())) // set to *type, not **type
					fieldVal = fieldVal.Elem()                   // use value, not *value
				} else {
					fieldVal.Set(reflect.New(field.Type).Elem())
				}

				if err := starlarkToGo(attrVal, fieldVal); err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported type: %s", srcType)
	}
}

func findStructFieldByTag(gotype reflect.Type, tagKey, tagValue string) (string, bool) {

	for i := 0; i < gotype.NumField(); i++ {
		field := gotype.Field(i)
		val, ok := field.Tag.Lookup(tagKey)
		if !ok {
			continue
		}
		if strings.EqualFold(val, tagValue) {
			return field.Name, true
		}
	}

	return "", false
}

func getExactMapType(val starlark.Value, gotype reflect.Type) reflect.Type {
	switch val.Type() {
	case "dict":
		if gotype.Kind() == reflect.Map {
			return gotype
		}
		return reflect.TypeOf(map[any]any{})
	default:
		return gotype
	}
}
