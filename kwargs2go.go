package startype

import (
	"fmt"
	"reflect"

	"go.starlark.net/starlark"
)

type KwargsValue struct {
	kwargs []starlark.Tuple
}

// Kwargs starts the conversion of a Starlark kwargs (keyword args) value
// to a Go struct. The function uses annotated fields on the struct to describe
// the keyword argument mapping as:
//
//	var Param struct {
//	    OutputFile  string   `name:"output_file" optional:"true"`
//	    SourcePaths []string `name:"output_path"`
//	}
//
// The struct can be mapped to the following keyword args:
//
//	kwargs := []starlark.Tuple{
//	   {starlark.String("output_file"), starlark.String("/tmp/out.tar.gz")},
//	   {starlark.String("source_paths"), starlark.NewList([]starlark.Value{starlark.String("/tmp/myfile")})},
//	}
//
// # Example
//
// Kwargs(kwargs).Go(&Param)
//
// Supported annotation: `name:"arg_name" optional:"true|false" (default false)`
func Kwargs(kwargs []starlark.Tuple) *KwargsValue {
	return &KwargsValue{kwargs: kwargs}
}

func (v *KwargsValue) Go(gostruct any) error {
	if v.kwargs == nil {
		return fmt.Errorf("keyword arguments is nil")
	}
	goval := reflect.ValueOf(gostruct)
	gotype := reflect.ValueOf(gostruct).Type()
	if gotype.Kind() != reflect.Pointer || goval.IsNil() {
		return fmt.Errorf("kwargs expects a non-nil pointer to a struct, got %v", gotype.Kind())
	}

	return kwargsToGo(v.kwargs, goval.Elem())
}

func kwargsToGo(kwargs []starlark.Tuple, goval reflect.Value) error {
	gotype := goval.Type()
	if gotype.Kind() != reflect.Struct {
		return fmt.Errorf("target type %s: a struct", gotype.Kind())
	}

	if !goval.IsValid() {
		goval.Set(reflect.Zero(goval.Type()))
	}

	for i := 0; i < goval.NumField(); i++ {
		field := gotype.Field(i)

		argName, ok := field.Tag.Lookup("name")
		if !ok {
			continue
		}

		// get arg from keyword args (use either tag or field name)
		kwarg, err := getKwarg(kwargs, argName, field.Name)
		if err != nil {
			return err
		}

		// is arg marked optional? By default args are optional=false
		// arg is optional if it is explicitly marked with "true" or "yes"
		argOptional, _ := field.Tag.Lookup("optional")
		switch argOptional {
		case "true", "yes":
		default:
			if kwarg == starlark.None {
				return fmt.Errorf("argument '%s' is required", argName)
			}
		}

		// set field value if not None
		if kwarg != starlark.None {
			fieldVal := goval.FieldByName(field.Name)

			if fieldVal.Kind() == reflect.Pointer {
				fieldVal.Set(reflect.New(field.Type.Elem()))
				fieldVal = fieldVal.Elem()
			} else {
				fieldVal.Set(reflect.New(field.Type).Elem())
			}

			if err := starlarkToGo(kwarg, fieldVal); err != nil {
				return err
			}
		}
	}

	return nil
}

func getKwarg(kwargs []starlark.Tuple, argName, defaultName string) (starlark.Value, error) {
	for _, kwarg := range kwargs {
		arg, ok := kwarg.Index(0).(starlark.String)
		if !ok {
			return nil, fmt.Errorf("keyword arg name is not a string")
		}
		if string(arg) == argName || string(arg) == defaultName {
			return kwarg.Index(1), nil
		}
	}
	return starlark.None, nil
}
