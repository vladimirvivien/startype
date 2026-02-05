package startype

import (
	"fmt"
	"reflect"
	"strconv"

	"go.starlark.net/starlark"
)

// ArgsValue holds both positional and keyword arguments for conversion
type ArgsValue struct {
	args   starlark.Tuple
	kwargs []starlark.Tuple
}

// Args creates a converter for both positional and keyword arguments.
// Positional args are matched by `position` struct tag.
// Keyword args are matched by `name` struct tag.
// If both provide a value for the same field, keyword wins.
//
// Example:
//
//	var params struct {
//	    Path     string `name:"path" position:"0" required:"true"`
//	    Encoding string `name:"encoding" position:"1"`
//	}
//	Args(args, kwargs).Go(&params)
func Args(args starlark.Tuple, kwargs []starlark.Tuple) *ArgsValue {
	return &ArgsValue{args: args, kwargs: kwargs}
}

// Go converts the arguments to a Go struct.
// The struct must use tags: `name`, `position`, `required`, `optional`
func (v *ArgsValue) Go(dest interface{}) error {
	destVal := reflect.ValueOf(dest)
	destType := destVal.Type()
	if destType.Kind() != reflect.Pointer || destVal.IsNil() {
		return fmt.Errorf("Args expects a non-nil pointer to a struct, got %v", destType.Kind())
	}
	return argsToGo(v.args, v.kwargs, destVal.Elem())
}

// fieldMeta holds metadata about a struct field for argument mapping
type fieldMeta struct {
	index    int
	position int    // -1 if not positional
	name     string // empty if no name tag
	required bool
}

func argsToGo(args starlark.Tuple, kwargs []starlark.Tuple, destVal reflect.Value) error {
	destType := destVal.Type()
	if destType.Kind() != reflect.Struct {
		return fmt.Errorf("destination must be a struct, got %s", destType.Kind())
	}

	// Build field metadata from struct tags
	fields := make([]fieldMeta, 0, destType.NumField())
	positionMap := make(map[int]int) // position -> fields index
	nameMap := make(map[string]int)  // name -> fields index

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		meta := fieldMeta{index: i, position: -1}

		// Get name tag (for kwargs matching)
		if name, ok := field.Tag.Lookup("name"); ok {
			meta.name = name
			nameMap[name] = len(fields)
		}

		// Get position tag (for positional args)
		if pos, ok := field.Tag.Lookup("position"); ok {
			if p, err := strconv.Atoi(pos); err == nil {
				meta.position = p
				positionMap[p] = len(fields)
			}
		}

		// Skip fields without name or position tags
		if meta.name == "" && meta.position < 0 {
			continue
		}

		// Get required tag
		if req, ok := field.Tag.Lookup("required"); ok {
			meta.required = req == "true" || req == "yes"
		}

		fields = append(fields, meta)
	}

	// Track which fields have been set
	setFields := make(map[int]bool)

	// 1. Process positional arguments first
	for i := 0; i < len(args); i++ {
		fieldIdx, ok := positionMap[i]
		if !ok {
			return fmt.Errorf("unexpected positional argument at index %d", i)
		}
		meta := fields[fieldIdx]
		fieldVal := destVal.Field(meta.index)

		if err := setFieldValue(fieldVal, args[i]); err != nil {
			return fmt.Errorf("positional arg %d: %w", i, err)
		}
		setFields[fieldIdx] = true
	}

	// 2. Process keyword arguments (can override positional)
	for _, kwarg := range kwargs {
		nameVal, ok := kwarg.Index(0).(starlark.String)
		if !ok {
			return fmt.Errorf("keyword argument name is not a string")
		}
		name := string(nameVal)

		fieldIdx, ok := nameMap[name]
		if !ok {
			return fmt.Errorf("unknown keyword argument: %s", name)
		}
		meta := fields[fieldIdx]
		fieldVal := destVal.Field(meta.index)

		if err := setFieldValue(fieldVal, kwarg.Index(1)); err != nil {
			return fmt.Errorf("keyword arg '%s': %w", name, err)
		}
		setFields[fieldIdx] = true
	}

	// 3. Validate required fields
	for i, meta := range fields {
		if meta.required && !setFields[i] {
			name := meta.name
			if name == "" {
				name = fmt.Sprintf("position %d", meta.position)
			}
			return fmt.Errorf("missing required argument: %s", name)
		}
	}

	return nil
}

// setFieldValue handles pointer allocation and calls starlarkToGo
func setFieldValue(fieldVal reflect.Value, val starlark.Value) error {
	if fieldVal.Kind() == reflect.Pointer {
		fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
		fieldVal = fieldVal.Elem()
	}
	return starlarkToGo(val, fieldVal)
}
