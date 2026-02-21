package startype

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// GoValue represents an inherent Go value which can be
// converted to a Starlark value/type
type GoValue[T any] struct {
	val T
}

// Go wraps a Go value into GoValue so that it can be converted to
// a Starlark value.
func Go[T any](val T) *GoValue[T] {
	return &GoValue[T]{val: val}
}

// Map wraps a Go map for clearer intent when converting to Starlark Dict.
func Map[K comparable, V any](m map[K]V) *GoValue[map[K]V] { return Go(m) }

// Slice wraps a Go slice for clearer intent when converting to Starlark List.
func Slice[T any](s []T) *GoValue[[]T] { return Go(s) }

// Value returns the original Go value with its concrete type.
func (v *GoValue[T]) Value() T {
	return v.val
}

// Starlark translates Go value to a starlark.Value value
// using the following type mapping:
//
//	bool                -- starlark.Bool
//	int{8,16,32,64}     -- starlark.Int
//	uint{8,16,32,64}    -- starlark.Int
//	float{32,64}        -- starlark.Float
//	string              -- starlark.String
//	[]T, [n]T           -- starlark.Tuple
//	map[K]T	            -- *starlark.Dict
//
// The specified Starlark value must be provided as
// a pointer to the target Starlark type.
//
// Example:
//
//	num := 64
//	var starInt starlark.Int
//	Go(num).Starlark(&starInt)
//
// For starlark.List and starlark.Set refer to their
// respective namesake methods.
func (v *GoValue[T]) Starlark(starval interface{}) error {
	return goToStarlark(v.val, starval)
}

// StarlarkList converts a slice of Go values to a starlark.Tuple,
// then converts that tuple into a starlark.List
func (v *GoValue[T]) StarlarkList(starval interface{}) error {
	var tuple starlark.Tuple
	if err := v.Starlark(&tuple); err != nil {
		return err
	}
	switch val := starval.(type) {
	case *starlark.Value:
		*val = starlark.NewList(tuple)
	case *starlark.List:
		*val = *starlark.NewList(tuple)
	case **starlark.List:
		listVal := *starlark.NewList(tuple)
		*val = &listVal
	default:
		return fmt.Errorf("target type %T: must be *starlark.List or *starlark.Value", starval)
	}
	return nil
}

// StarlarkSet converts a slice of Go values to a starlark.Tuple,
// then converts that tuple into a starlark.Set
func (v *GoValue[T]) StarlarkSet(starval interface{}) error {
	var tuple starlark.Tuple
	if err := v.Starlark(&tuple); err != nil {
		return err
	}

	starSet := starlark.NewSet(len(tuple))
	for _, val := range tuple {
		if err := starSet.Insert(val); err != nil {
			continue
		}
	}

	switch val := starval.(type) {
	case *starlark.Value:
		*val = starSet
	case **starlark.Set:
		*val = starSet
	default:
		return fmt.Errorf("target type %T: must be **starlark.Set or *starlark.Value", starval)
	}
	return nil
}

// GoStructToStringDict is a helper func that converts a Go struct type to
// starlark.StringDict.
func GoStructToStringDict(gostruct interface{}) (starlark.StringDict, error) {
	goval := reflect.ValueOf(gostruct)
	gotype := goval.Type()
	if gotype.Kind() != reflect.Struct {
		return nil, fmt.Errorf("source type must be a struct")
	}
	return goStructToStringDict(goval)
}

// goToStarlark translates Go value to a starlark.Value value
// using the following type mapping:
//
//		bool				-- starlark.Bool
//		int{8,16,32,64}		-- starlark.Int
//		uint{8,16,32,64}	-- starlark.Int
//		float{32,64}		-- starlark.Float
//	    string			 	-- starlark.String
//	    []T, [n]T			-- starlark.Tuple
//		map[K]T				-- *starlark.Dict
func goToStarlark(gov interface{}, starval interface{}) error {
	if gov == nil {
		if val, ok := starval.(*starlark.Value); ok {
			*val = starlark.None
		}
		return nil
	}
	goval := reflect.ValueOf(gov)
	if !goval.IsValid() {
		if val, ok := starval.(*starlark.Value); ok {
			*val = starlark.None
		}
		return nil
	}

	gotype := goval.Type()
	switch gotype.Kind() {
	case reflect.Bool:
		switch val := starval.(type) {
		case *starlark.Value:
			*val = starlark.Bool(goval.Bool())
		case *starlark.Bool:
			*val = starlark.Bool(goval.Bool())
		case **starlark.Bool:
			boolVal := starlark.Bool(goval.Bool())
			*val = &boolVal
		default:
			return fmt.Errorf("target type (%T): must be *starlark.Bool, *starlark.Value", starval)
		}

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch val := starval.(type) {
		case *starlark.Value:
			*val = starlark.MakeInt64(goval.Int())
		case *starlark.Int:
			*val = starlark.MakeInt64(goval.Int())
		case **starlark.Int:
			intVal := starlark.MakeInt64(goval.Int())
			*val = &intVal
		default:
			return fmt.Errorf("target type (%T): must be *starlark.Int or *starlark.Value", starval)
		}
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch val := starval.(type) {
		case *starlark.Value:
			*val = starlark.MakeUint64(goval.Uint())
		case *starlark.Int:
			*val = starlark.MakeUint64(goval.Uint())
		case **starlark.Int:
			uintVal := starlark.MakeUint64(goval.Uint())
			*val = &uintVal
		default:
			return fmt.Errorf("target type %T: must be *starlark.Int or *starlark.Value", starval)
		}
		return nil

	case reflect.Float32, reflect.Float64:
		switch val := starval.(type) {
		case *starlark.Value:
			*val = starlark.Float(goval.Float())
		case *starlark.Float:
			*val = starlark.Float(goval.Float())
		case **starlark.Float:
			floatVal := starlark.Float(goval.Float())
			*val = &floatVal
		default:
			return fmt.Errorf("target type %T: must be *starlark.Float or *starlark.Value", starval)
		}
		return nil

	case reflect.String:
		switch val := starval.(type) {
		case *starlark.Value:
			*val = starlark.String(goval.String())
		case *starlark.String:
			*val = starlark.String(goval.String())
		case **starlark.String:
			strVal := starlark.String(goval.String())
			*val = &strVal
		default:
			return fmt.Errorf("target type %T: must be *starlark.String or *starlark.Value", starval)
		}
		return nil

	case reflect.Slice, reflect.Array:
		result, err := makeTuple(goval)
		if err != nil {
			return err
		}

		switch val := starval.(type) {
		case *starlark.Value:
			*val = starlark.NewList(result)
		case *starlark.Tuple:
			*val = result
		case **starlark.Tuple:
			tupVal := starlark.Tuple(result)
			*val = &tupVal
		case *starlark.List:
			*val = *starlark.NewList(result)
		case **starlark.List:
			listVal := starlark.NewList(result)
			*val = listVal
		default:
			return fmt.Errorf("target type %T: must be *starlark.Tuple, *starlark.List, or *starlark.Value", starval)
		}

		return nil

	case reflect.Map:
		dict, err := goMapToDict(goval)
		if err != nil {
			return err
		}

		switch dictVal := starval.(type) {
		case **starlark.Dict:
			*dictVal = dict
		case *starlark.Value:
			*dictVal = dict
		default:
			return fmt.Errorf("want target type **starlark.Dict or *starlark.Value, got: %T", dictVal)
		}

		return nil

	case reflect.Struct:
		dict, err := goStructToStringDict(goval)
		if err != nil {
			return err
		}

		switch val := starval.(type) {
		case *starlark.Value:
			result := starlarkstruct.FromStringDict(starlark.String(gotype.Name()), dict)
			*val = result
		case *starlarkstruct.Struct:
			result := starlarkstruct.FromStringDict(starlark.String(gotype.Name()), dict)
			*val = *result
		case **starlarkstruct.Struct:
			result := starlarkstruct.FromStringDict(starlark.String(gotype.Name()), dict)
			*val = result
		case *starlark.StringDict:
			*val = dict
		case **starlark.StringDict:
			*val = &dict
		default:
			return fmt.Errorf("target type %T: must be *starlarkstruct.Struct or *starlark.Value", starval)
		}

		return nil

	case reflect.Pointer:
		goElem := goval.Elem()
		if !goElem.IsValid() {
			return nil
		}
		return goToStarlark(goElem.Interface(), starval)

	default:
		return fmt.Errorf("unable to convert Go type %T to Starlark type", gov)
	}

}

func makeTuple(sliceVal reflect.Value) ([]starlark.Value, error) {
	tuple := make([]starlark.Value, sliceVal.Len())
	for i := 0; i < sliceVal.Len(); i++ {
		var elem starlark.Value
		if err := goToStarlark(sliceVal.Index(i).Interface(), &elem); err != nil {
			return nil, err
		}
		tuple[i] = elem
	}
	return tuple, nil
}

func goMapToDict(mapVal reflect.Value) (*starlark.Dict, error) {
	iter := mapVal.MapRange()
	dict := starlark.NewDict(mapVal.Len())

	for iter.Next() {
		// convert key
		var key starlark.Value
		if err := goToStarlark(iter.Key().Interface(), &key); err != nil {
			return nil, fmt.Errorf("GoToStarlrk: failed map key conversion: %s", err)
		}

		// convert value
		var val starlark.Value
		if err := goToStarlark(iter.Value().Interface(), &val); err != nil {
			return nil, fmt.Errorf("GoToStarlark: failed map value conversion: %s", err)
		}

		if err := dict.SetKey(key, val); err != nil {
			return nil, fmt.Errorf("GoToStarlark: failed to set map value with key: %s", key)
		}
	}
	return dict, nil
}

func goStructToStringDict(goval reflect.Value) (starlark.StringDict, error) {
	gotype := goval.Type()
	stringDict := make(starlark.StringDict)
	for i := 0; i < goval.NumField(); i++ {
		field := gotype.Field(i)
		// only grab exported field to avoid panic
		if !field.IsExported() {
			continue
		}

		fname := field.Name
		// get starlarkstruct field name from tag (if any)
		name, _ := field.Tag.Lookup("name")
		if name != "" {
			fname = name
		}

		var fval starlark.Value

		if err := goToStarlark(goval.Field(i).Interface(), &fval); err != nil {
			return nil, fmt.Errorf("GoToStarlark: failed struct field conversion: %s", err)
		}
		stringDict[fname] = fval
	}

	return stringDict, nil
}

// --- Dynamic dispatch: any → starlark.Value ---

// ToStarlarkValue performs dynamic dispatch to convert the wrapped Go value
// to a starlark.Value. It handles: nil→None, bool→Bool, int types→Int,
// float64→Int|Float (JSON semantics: integer floats→Int), string→String,
// []any→List (recursive), map[string]any→Dict (sorted keys, recursive).
// For other slice/map types, it falls back to reflect-based iteration.
func (v *GoValue[T]) ToStarlarkValue() (starlark.Value, error) {
	return anyToStarlarkValue(v.val)
}

// ToBool converts the wrapped Go value to starlark.Bool.
func (v *GoValue[T]) ToBool() (starlark.Bool, error) {
	if b, ok := any(v.val).(bool); ok {
		return starlark.Bool(b), nil
	}
	return starlark.False, fmt.Errorf("ToBool: value is %T, not bool", v.val)
}

// ToInt converts the wrapped Go integer value to starlark.Int.
func (v *GoValue[T]) ToInt() (starlark.Int, error) {
	switch val := any(v.val).(type) {
	case int:
		return starlark.MakeInt(val), nil
	case int8:
		return starlark.MakeInt64(int64(val)), nil
	case int16:
		return starlark.MakeInt64(int64(val)), nil
	case int32:
		return starlark.MakeInt64(int64(val)), nil
	case int64:
		return starlark.MakeInt64(val), nil
	case uint:
		return starlark.MakeUint64(uint64(val)), nil
	case uint8:
		return starlark.MakeUint64(uint64(val)), nil
	case uint16:
		return starlark.MakeUint64(uint64(val)), nil
	case uint32:
		return starlark.MakeUint64(uint64(val)), nil
	case uint64:
		return starlark.MakeUint64(val), nil
	}
	return starlark.Int{}, fmt.Errorf("ToInt: value is %T, not an integer type", v.val)
}

// ToFloat converts the wrapped Go float value to starlark.Float.
func (v *GoValue[T]) ToFloat() (starlark.Float, error) {
	switch val := any(v.val).(type) {
	case float32:
		return starlark.Float(float64(val)), nil
	case float64:
		return starlark.Float(val), nil
	}
	return 0, fmt.Errorf("ToFloat: value is %T, not a float type", v.val)
}

// ToString converts the wrapped Go string value to starlark.String.
func (v *GoValue[T]) ToString() (starlark.String, error) {
	if s, ok := any(v.val).(string); ok {
		return starlark.String(s), nil
	}
	return "", fmt.Errorf("ToString: value is %T, not string", v.val)
}

// ToDict converts the wrapped Go map to a *starlark.Dict.
// Keys are sorted for deterministic output.
func (v *GoValue[T]) ToDict() (*starlark.Dict, error) {
	rv := reflect.ValueOf(v.val)
	if !rv.IsValid() || rv.Kind() != reflect.Map {
		return nil, fmt.Errorf("ToDict: value is %T, not a map", v.val)
	}
	return reflectMapToDict(rv)
}

// ToList converts the wrapped Go slice/array to a *starlark.List.
func (v *GoValue[T]) ToList() (*starlark.List, error) {
	rv := reflect.ValueOf(v.val)
	if !rv.IsValid() || (rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array) {
		return nil, fmt.Errorf("ToList: value is %T, not a slice or array", v.val)
	}
	return reflectSliceToList(rv)
}

// anyToStarlarkValue converts an arbitrary Go value to a starlark.Value
// using dynamic type dispatch. This is the core implementation for
// ToStarlarkValue and is also used by container converters recursively.
func anyToStarlarkValue(v any) (starlark.Value, error) {
	switch val := v.(type) {
	case nil:
		return starlark.None, nil
	case bool:
		return starlark.Bool(val), nil
	case int:
		return starlark.MakeInt(val), nil
	case int8:
		return starlark.MakeInt64(int64(val)), nil
	case int16:
		return starlark.MakeInt64(int64(val)), nil
	case int32:
		return starlark.MakeInt64(int64(val)), nil
	case int64:
		return starlark.MakeInt64(val), nil
	case uint:
		return starlark.MakeUint64(uint64(val)), nil
	case uint8:
		return starlark.MakeUint64(uint64(val)), nil
	case uint16:
		return starlark.MakeUint64(uint64(val)), nil
	case uint32:
		return starlark.MakeUint64(uint64(val)), nil
	case uint64:
		return starlark.MakeUint64(val), nil
	case float32:
		return starlark.Float(float64(val)), nil
	case float64:
		// JSON number semantics: integer floats become starlark.Int
		if val == math.Trunc(val) && !math.IsInf(val, 0) && !math.IsNaN(val) {
			return starlark.MakeInt64(int64(val)), nil
		}
		return starlark.Float(val), nil
	case string:
		return starlark.String(val), nil
	case []any:
		elems := make([]starlark.Value, len(val))
		for i, elem := range val {
			sv, err := anyToStarlarkValue(elem)
			if err != nil {
				return nil, fmt.Errorf("list[%d]: %w", i, err)
			}
			elems[i] = sv
		}
		return starlark.NewList(elems), nil
	case map[string]any:
		dict := starlark.NewDict(len(val))
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sv, err := anyToStarlarkValue(val[k])
			if err != nil {
				return nil, fmt.Errorf("dict[%q]: %w", k, err)
			}
			if err := dict.SetKey(starlark.String(k), sv); err != nil {
				return nil, fmt.Errorf("dict set %q: %w", k, err)
			}
		}
		return dict, nil
	case starlark.Value:
		return val, nil
	default:
		// Fall back to reflect for other slice/map types
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			list, err := reflectSliceToList(rv)
			if err != nil {
				return nil, err
			}
			return list, nil
		case reflect.Map:
			dict, err := reflectMapToDict(rv)
			if err != nil {
				return nil, err
			}
			return dict, nil
		}
		return nil, fmt.Errorf("unsupported Go type %T for dynamic conversion", v)
	}
}

// reflectSliceToList converts a reflect.Value slice/array to *starlark.List.
func reflectSliceToList(rv reflect.Value) (*starlark.List, error) {
	elems := make([]starlark.Value, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		sv, err := anyToStarlarkValue(rv.Index(i).Interface())
		if err != nil {
			return nil, fmt.Errorf("list[%d]: %w", i, err)
		}
		elems[i] = sv
	}
	return starlark.NewList(elems), nil
}

// reflectMapToDict converts a reflect.Value map to *starlark.Dict with sorted keys.
func reflectMapToDict(rv reflect.Value) (*starlark.Dict, error) {
	dict := starlark.NewDict(rv.Len())

	// Collect and sort keys for deterministic output
	keys := rv.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return fmt.Sprint(keys[i].Interface()) < fmt.Sprint(keys[j].Interface())
	})

	for _, k := range keys {
		key, err := anyToStarlarkValue(k.Interface())
		if err != nil {
			return nil, fmt.Errorf("dict key: %w", err)
		}
		val, err := anyToStarlarkValue(rv.MapIndex(k).Interface())
		if err != nil {
			return nil, fmt.Errorf("dict value: %w", err)
		}
		if err := dict.SetKey(key, val); err != nil {
			return nil, err
		}
	}
	return dict, nil
}
