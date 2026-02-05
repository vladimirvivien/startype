package startype

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// GoValue represents an inherent Go value which can be
// converted to a Starlark value/type
type GoValue struct {
	val interface{}
}

// Go wraps a Go value into GoValue so that it can be converted to
// a Starlark value.
func Go(val interface{}) *GoValue {
	return &GoValue{val: val}
}

// Value returns the original Go value as an interface{}
func (v *GoValue) Value() interface{} {
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
func (v *GoValue) Starlark(starval interface{}) error {
	return goToStarlark(v.val, starval)
}

// StarlarkList converts a slice of Go values to a starlark.Tuple,
// then converts that tuple into a starlark.List
func (v *GoValue) StarlarkList(starval interface{}) error {
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
func (v *GoValue) StarlarkSet(starval interface{}) error {
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
	case *starlark.Set:
		*val = *starSet //nolint:govet
	case **starlark.Set:
		*val = starSet
	default:
		return fmt.Errorf("target type %T: must be *starlark.List or *starlark.Value", starval)
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
		case *starlark.Dict:
			*dictVal = *dict //nolint:govet
		case **starlark.Dict:
			*dictVal = dict
		case *starlark.Value:
			*dictVal = dict
		default:
			return fmt.Errorf("want target type *starlark.Dict, got: %T", dictVal)
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
