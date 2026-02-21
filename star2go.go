package startype

import (
	"fmt"
	"reflect"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// StarValue represents a wrapped Starlark value which can be
// converted to a Go value.
type StarValue[T starlark.Value] struct {
	val T
}

// Starlark wraps a Starlark value val
// so it can be converted to a Go value.
func Starlark[T starlark.Value](val T) *StarValue[T] {
	return &StarValue[T]{val: val}
}

// Dict wraps a *starlark.Dict for clearer intent when converting to Go map.
func Dict(d *starlark.Dict) *StarValue[*starlark.Dict] { return Starlark(d) }

// List wraps a *starlark.List for clearer intent when converting to Go slice.
func List(l *starlark.List) *StarValue[*starlark.List] { return Starlark(l) }

// Value returns the wrapped Starlark value with its concrete type.
func (v *StarValue[T]) Value() T {
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

func (v *StarValue[T]) Go(goin interface{}) error {
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

	// Handle passthrough types - assign directly without conversion
	// Note: Check Callable before Value since Callable embeds Value

	// starlark.Callable - accept callable values
	starlarkCallableType := reflect.TypeOf((*starlark.Callable)(nil)).Elem()
	if gotype == starlarkCallableType {
		if callable, ok := srcVal.(starlark.Callable); ok {
			goval.Set(reflect.ValueOf(callable))
			return nil
		}
		return fmt.Errorf("value is not callable: got %s", srcVal.Type())
	}

	// starlark.Value - accept any Starlark value
	starlarkValueType := reflect.TypeOf((*starlark.Value)(nil)).Elem()
	if gotype == starlarkValueType {
		goval.Set(reflect.ValueOf(srcVal))
		return nil
	}

	// starlark.Bytes - pass through if target is starlark.Bytes
	starlarkBytesType := reflect.TypeOf(starlark.Bytes(""))
	if gotype == starlarkBytesType {
		if bytes, ok := srcVal.(starlark.Bytes); ok {
			goval.Set(reflect.ValueOf(bytes))
			return nil
		}
		return fmt.Errorf("value is not bytes: got %s", srcVal.Type())
	}

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
		listVal, ok := srcVal.(*starlark.List)
		if !ok {
			return fmt.Errorf("failed to assert %T as *starlark.List", srcVal)
		}
		switch gotype.Kind() {
		case reflect.Slice, reflect.Array:
			goval.Set(reflect.MakeSlice(gotype, listVal.Len(), listVal.Len()))
			for i := 0; i < listVal.Len(); i++ {
				if err := starlarkToGo(listVal.Index(i), goval.Index(i)); err != nil {
					return err
				}
			}
		case reflect.Interface:
			result := make([]any, listVal.Len())
			for i := 0; i < listVal.Len(); i++ {
				elem := reflect.New(reflect.TypeOf((*any)(nil)).Elem()).Elem()
				if err := starlarkToGo(listVal.Index(i), elem); err != nil {
					return err
				}
				result[i] = elem.Interface()
			}
			goval.Set(reflect.ValueOf(result))
		default:
			return fmt.Errorf("target type must be slice, array, or any")
		}
		return nil

	case "tuple":
		tupVal, ok := srcVal.(starlark.Tuple)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Tuple", srcVal)
		}
		switch gotype.Kind() {
		case reflect.Slice, reflect.Array:
			goval.Set(reflect.MakeSlice(gotype, tupVal.Len(), tupVal.Len()))
			for i := 0; i < tupVal.Len(); i++ {
				if err := starlarkToGo(tupVal.Index(i), goval.Index(i)); err != nil {
					return err
				}
			}
		case reflect.Interface:
			result := make([]any, tupVal.Len())
			for i := 0; i < tupVal.Len(); i++ {
				elem := reflect.New(reflect.TypeOf((*any)(nil)).Elem()).Elem()
				if err := starlarkToGo(tupVal.Index(i), elem); err != nil {
					return err
				}
				result[i] = elem.Interface()
			}
			goval.Set(reflect.ValueOf(result))
		default:
			return fmt.Errorf("target type must be slice, array, or any")
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

		// map target type — when target is interface{}, create a concrete map
		// and use mapVal to track the actual map for SetMapIndex calls
		var mapVal reflect.Value
		switch gotype.Kind() {
		case reflect.Map:
			mapVal = reflect.MakeMapWithSize(gotype, dict.Len())
			goval.Set(mapVal)
		case reflect.Interface:
			gotype = reflect.TypeOf(map[any]any{})
			mapVal = reflect.MakeMapWithSize(gotype, dict.Len())
			goval.Set(mapVal)
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
			mapVal.SetMapIndex(goMapKey, goMapElem)
		}
		return nil

	case "set":
		setVal, ok := srcVal.(*starlark.Set)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Set", srcVal)
		}
		switch gotype.Kind() {
		case reflect.Slice, reflect.Array:
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
		case reflect.Interface:
			result := make([]any, setVal.Len())
			var setItem starlark.Value
			iter := setVal.Iterate()
			i := 0
			for iter.Next(&setItem) {
				elem := reflect.New(reflect.TypeOf((*any)(nil)).Elem()).Elem()
				if err := starlarkToGo(setItem, elem); err != nil {
					return err
				}
				result[i] = elem.Interface()
				i++
			}
			goval.Set(reflect.ValueOf(result))
		default:
			return fmt.Errorf("target type must be slice, array, or any")
		}
		return nil

	case "bytes":
		bytesVal, ok := srcVal.(starlark.Bytes)
		if !ok {
			return fmt.Errorf("failed to assert %T as starlark.Bytes", srcVal)
		}

		switch gotype.Kind() {
		case reflect.Slice:
			if gotype.Elem().Kind() == reflect.Uint8 {
				// Convert starlark.Bytes → []byte
				goval.Set(reflect.ValueOf([]byte(bytesVal)))
				return nil
			}
			return fmt.Errorf("starlark.Bytes to Go: slice element must be uint8, got %s", gotype.Elem().Kind())
		case reflect.Interface:
			// For `any` target, convert to []byte
			goval.Set(reflect.ValueOf([]byte(bytesVal)))
			return nil
		default:
			return fmt.Errorf("starlark.Bytes to Go: target type (%s) must be []byte or any", gotype.Kind())
		}

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
	case "NoneType":
		if gotype.Kind() == reflect.Interface {
			// leave as zero value (nil) for interface targets
			return nil
		}
		return fmt.Errorf("NoneType: target type (%s) must be any", gotype.Kind())

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

// --- Dynamic dispatch: starlark.Value → any ---

// ToGoValue performs dynamic dispatch to convert the wrapped Starlark value
// to a Go value. It handles: None→nil, Bool→bool, Int→int64, Float→float64,
// String→string, List→[]any (recursive), Tuple→[]any, Dict→map[string]any
// (recursive, requires string keys). Unknown types fall back to String().
func (v *StarValue[T]) ToGoValue() (any, error) {
	return starlarkValueToGo(v.val)
}

// ToBool converts the wrapped Starlark value to a Go bool.
func (v *StarValue[T]) ToBool() (bool, error) {
	if b, ok := any(v.val).(starlark.Bool); ok {
		return bool(b), nil
	}
	return false, fmt.Errorf("ToBool: value is %s, not bool", any(v.val).(starlark.Value).Type())
}

// ToInt64 converts the wrapped Starlark value to a Go int64.
func (v *StarValue[T]) ToInt64() (int64, error) {
	if i, ok := any(v.val).(starlark.Int); ok {
		val, ok := i.Int64()
		if !ok {
			return 0, fmt.Errorf("ToInt64: value too large for int64")
		}
		return val, nil
	}
	return 0, fmt.Errorf("ToInt64: value is %s, not int", any(v.val).(starlark.Value).Type())
}

// ToFloat64 converts the wrapped Starlark value to a Go float64.
func (v *StarValue[T]) ToFloat64() (float64, error) {
	if f, ok := any(v.val).(starlark.Float); ok {
		return float64(f), nil
	}
	return 0, fmt.Errorf("ToFloat64: value is %s, not float", any(v.val).(starlark.Value).Type())
}

// ToString converts the wrapped Starlark value to a Go string.
func (v *StarValue[T]) ToString() (string, error) {
	if s, ok := any(v.val).(starlark.String); ok {
		return string(s), nil
	}
	return "", fmt.Errorf("ToString: value is %s, not string", any(v.val).(starlark.Value).Type())
}

// ToMap converts the wrapped Starlark Dict to a Go map[string]any.
func (v *StarValue[T]) ToMap() (map[string]any, error) {
	dict, ok := any(v.val).(*starlark.Dict)
	if !ok {
		return nil, fmt.Errorf("ToMap: value is %s, not dict", any(v.val).(starlark.Value).Type())
	}
	result := make(map[string]any, dict.Len())
	for _, kv := range dict.Items() {
		key, ok := kv[0].(starlark.String)
		if !ok {
			return nil, fmt.Errorf("ToMap: dict key must be string, got %s", kv[0].Type())
		}
		val, err := starlarkValueToGo(kv[1])
		if err != nil {
			return nil, fmt.Errorf("ToMap: dict[%q]: %w", string(key), err)
		}
		result[string(key)] = val
	}
	return result, nil
}

// ToSlice converts the wrapped Starlark List to a Go []any.
func (v *StarValue[T]) ToSlice() ([]any, error) {
	list, ok := any(v.val).(*starlark.List)
	if !ok {
		return nil, fmt.Errorf("ToSlice: value is %s, not list", any(v.val).(starlark.Value).Type())
	}
	result := make([]any, list.Len())
	for i := 0; i < list.Len(); i++ {
		val, err := starlarkValueToGo(list.Index(i))
		if err != nil {
			return nil, fmt.Errorf("ToSlice: list[%d]: %w", i, err)
		}
		result[i] = val
	}
	return result, nil
}

// starlarkValueToGo converts any starlark.Value to a Go value using
// dynamic type dispatch. This is the core implementation shared by
// ToGoValue, ToMap, and ToSlice.
func starlarkValueToGo(v starlark.Value) (any, error) {
	switch val := v.(type) {
	case starlark.NoneType:
		return nil, nil
	case starlark.Bool:
		return bool(val), nil
	case starlark.Int:
		if i, ok := val.Int64(); ok {
			return i, nil
		}
		// Fall back to string for very large integers
		return val.BigInt().String(), nil
	case starlark.Float:
		return float64(val), nil
	case starlark.String:
		return string(val), nil
	case *starlark.List:
		result := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			item, err := starlarkValueToGo(val.Index(i))
			if err != nil {
				return nil, fmt.Errorf("list[%d]: %w", i, err)
			}
			result[i] = item
		}
		return result, nil
	case starlark.Tuple:
		result := make([]any, len(val))
		for i, item := range val {
			v, err := starlarkValueToGo(item)
			if err != nil {
				return nil, fmt.Errorf("tuple[%d]: %w", i, err)
			}
			result[i] = v
		}
		return result, nil
	case *starlark.Dict:
		result := make(map[string]any, val.Len())
		for _, kv := range val.Items() {
			key, ok := kv[0].(starlark.String)
			if !ok {
				return nil, fmt.Errorf("dict key must be string, got %s", kv[0].Type())
			}
			v, err := starlarkValueToGo(kv[1])
			if err != nil {
				return nil, fmt.Errorf("dict[%q]: %w", string(key), err)
			}
			result[string(key)] = v
		}
		return result, nil
	default:
		// Fall back to String() representation for unknown types
		return v.String(), nil
	}
}
