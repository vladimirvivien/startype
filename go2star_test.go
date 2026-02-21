package startype

import (
	"math"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestGoToStarlark(t *testing.T) {
	tests := []struct {
		name  string
		goVal interface{}
		eval  func(*testing.T, interface{})
	}{
		{
			name:  "bool",
			goVal: true,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Bool
				if err := Go(true).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval != true {
					t.Errorf("unexpected bool value: %t", starval)
				}
			},
		},
		{
			name:  "bool-pointer",
			goVal: true,
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Bool
				if err := Go(true).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if *starval != true {
					t.Errorf("unexpected bool value: %t", *starval)
				}
			},
		},
		{
			name:  "bool-value",
			goVal: true,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(true).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Truth() != true {
					t.Errorf("unexpected bool value: %t", starval)
				}
			},
		},
		{
			name:  "int",
			goVal: math.MaxInt32,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Int
				if err := Go(math.MaxInt32).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.Int64()
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val != math.MaxInt32 {
					t.Errorf("unexpected Int64 value: %d", val)
				}
			},
		},
		{
			name:  "int-value",
			goVal: math.MaxInt32,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(math.MaxInt32).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.(starlark.Int)
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val.BigInt().Int64() != math.MaxInt32 {
					t.Errorf("unexpected Int64 value: %d", val)
				}
			},
		},
		{
			name:  "int-pointer",
			goVal: math.MaxInt32,
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Int
				if err := Go(math.MaxInt32).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.Int64()
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val != math.MaxInt32 {
					t.Errorf("unexpected Int64 value: %d", val)
				}
			},
		},
		{
			name:  "uint",
			goVal: uint64(math.MaxUint64),
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Int
				if err := Go(uint64(math.MaxUint64)).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.Uint64()
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val != math.MaxUint64 {
					t.Errorf("unexpected Uint64 value: %d", val)
				}
			},
		},
		{
			name:  "uint-pointer",
			goVal: uint64(math.MaxUint64),
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Int
				if err := Go(uint64(math.MaxUint64)).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.Uint64()
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val != math.MaxUint64 {
					t.Errorf("unexpected Uint64 value: %d", val)
				}
			},
		},
		{
			name:  "uint-value",
			goVal: uint64(math.MaxUint64),
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(uint64(math.MaxUint64)).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval.(starlark.Int)
				if !ok {
					t.Errorf("starlark.Int.Int64 failed")
				}
				if val.BigInt().Uint64() != math.MaxUint64 {
					t.Errorf("unexpected Uint64 value: %d", val)
				}
			},
		},
		{
			name:  "float",
			goVal: math.MaxFloat32,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Float
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval != math.MaxFloat32 {
					t.Errorf("unexpected float value: %s", starval)
				}
			},
		},
		{
			name:  "float-pointer",
			goVal: math.MaxFloat32,
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Float
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if *starval != math.MaxFloat32 {
					t.Errorf("unexpected float-pointer value: %s", starval)
				}
			},
		},
		{
			name:  "float-value",
			goVal: math.MaxFloat32,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.(starlark.Float) != math.MaxFloat32 {
					t.Errorf("unexpected float value: %s", starval)
				}
			},
		},
		{
			name:  "string",
			goVal: "Hello World!",
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.String
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if string(starval) != `Hello World!` {
					t.Errorf("unexpected string value: %s", starval)
				}
			},
		},
		{
			name:  "string-pointer",
			goVal: "Hello World!",
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.String
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if string(*starval) != `Hello World!` {
					t.Errorf("unexpected string value: %s", starval)
				}
			},
		},
		{
			name:  "string-value",
			goVal: "Hello World!",
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.String() != `"Hello World!"` {
					t.Errorf("unexpected string value: %s", starval)
				}
			},
		},
		{
			name:  "tuple-string",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				starval := make(starlark.Tuple, 2)
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}
				if starval.Index(1).String() != `"World!"` {
					t.Errorf("unexpected value: %s", starval.Index(1).String())
				}
			},
		},
		{
			name:  "tuple-numeric",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				starval := make(starlark.Tuple, 3)
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 3 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}

				intVal, _ := starval.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "tuple-mix",
			goVal: []interface{}{1, 2, 3, "Go!"},
			eval: func(t *testing.T, goVal interface{}) {
				starval := make(starlark.Tuple, 4)
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 4 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}
				strVal := starval.Index(3).String()
				if strVal != `"Go!"` {
					t.Errorf("Unexpected string element: %s", strVal)
				}
			},
		},
		{
			name:  "tuple-value",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				list := starval.(*starlark.List)
				if list.Len() != 3 {
					t.Errorf("unexpected list length %d", list.Len())
				}

				intVal, _ := list.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "tuple-pointer",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Tuple
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if len(*starval) != 3 {
					t.Errorf("unexpected tuple length %d", len(*starval))
				}

				intVal, _ := starval.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "dict-strig-string",
			goVal: map[string]string{"msg": "hello", "target": "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Dict
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected dict length %d", starval.Len())
				}
				val, _, err := starval.Get(starlark.String("target"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}

				if val.String() != `"world"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "dict-string-int",
			goVal: map[string]int{"one": 12, "two": math.MaxInt8, "three": math.MaxInt64},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Dict
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 3 {
					t.Errorf("unexpected dict length %d", starval.Len())
				}
				val, _, err := starval.Get(starlark.String("three"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}
				if intVal, _ := val.(starlark.Int).Int64(); intVal != math.MaxInt64 {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "dict-any-mix",
			goVal: map[any]any{"one": "12", "two": math.MaxInt8, 3: math.MaxInt64},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Dict
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 3 {
					t.Errorf("unexpected dict length %d", starval.Len())
				}
				val, _, err := starval.Get(starlark.String("one"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}
				if strVal, _ := val.(starlark.String); strVal != "12" {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

				val, _, err = starval.Get(starlark.MakeInt(3))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}
				if intVal, _ := val.(starlark.Int).Int64(); intVal != math.MaxInt64 {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "dict-pointer",
			goVal: map[string]string{"msg": "hello", "target": "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Dict
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected dict length %d", starval.Len())
				}
				val, _, err := starval.Get(starlark.String("target"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}

				if val.String() != `"world"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "dict-value",
			goVal: map[string]int{"one": 12, "two": math.MaxInt8, "three": math.MaxInt64},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				dictVal := starval.(*starlark.Dict)
				if dictVal.Len() != 3 {
					t.Errorf("unexpected dict length %d", dictVal.Len())
				}
				val, _, err := dictVal.Get(starlark.String("three"))
				if err != nil {
					t.Errorf("failed to get value from starlark.Dict: %s", err)
				}
				if intVal, _ := val.(starlark.Int).Int64(); intVal != math.MaxInt64 {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "struct",
			goVal: struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlarkstruct.Struct
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, err := starval.Attr("Msg")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "struct-pointer",
			goVal: struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlarkstruct.Struct
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, err := starval.Attr("Msg")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "struct-stringdict",
			goVal: struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.StringDict
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, ok := starval["Msg"]
				if !ok {
					t.Fatalf("failed to get value from struct value from starlark.StringDict['Msg']")
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "struct-stringdict-pointer",
			goVal: struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.StringDict
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, ok := (*starval)["Msg"]
				if !ok {
					t.Fatalf("failed to get value from struct value from starlark.StringDict['Msg']")
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "struct-value",
			goVal: struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				structVal := starval.(*starlarkstruct.Struct)
				val, err := structVal.Attr("Msg")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name: "struct-annotated",
			goVal: struct {
				Msg    string `name:"msg_field"`
				Target string
			}{
				Msg: "hello", Target: "world",
			},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlarkstruct.Struct
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, err := starval.Attr("msg_field")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name: "struct-value-annotated",
			goVal: struct {
				Msg    string
				Target string `name:"tgt"`
			}{
				Msg: "hello", Target: "world",
			},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				structVal := starval.(*starlarkstruct.Struct)
				val, err := structVal.Attr("tgt")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				if val.String() != `"world"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}
			},
		},
		{
			name:  "struct-embedded-struct",
			goVal: struct{ Msg struct{ Message string } }{Msg: struct{ Message string }{Message: "Hello World!"}},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlarkstruct.Struct
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				structVal, err := starval.Attr("Msg")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: %s", err)
				}

				starstruct, ok := structVal.(*starlarkstruct.Struct)
				if !ok {
					t.Fatalf("unexpected type: %T", structVal)
				}

				val, err := starstruct.Attr("Message")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: [%#v], %s ", starval, err)
				}

				if val.String() != `"Hello World!"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}
			},
		},
		{
			name:  "nil-to-starlark-value",
			goVal: nil,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go[any](nil).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval != starlark.None {
					t.Errorf("expected starlark.None, got %v (%s)", starval, starval.Type())
				}
			},
		},
		{
			name:  "slice-to-starlark-value",
			goVal: []string{"a", "b", "c"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(goVal).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Type() != "list" {
					t.Errorf("expected list type, got %s", starval.Type())
				}
				list := starval.(*starlark.List)
				if list.Len() != 3 {
					t.Errorf("expected 3 elements, got %d", list.Len())
				}
			},
		},
		{
			name:  "zero-int-to-starlark",
			goVal: 0,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(0).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				intVal, ok := starval.(starlark.Int)
				if !ok {
					t.Fatalf("expected starlark.Int, got %T", starval)
				}
				v, _ := intVal.Int64()
				if v != 0 {
					t.Errorf("expected 0, got %d", v)
				}
			},
		},
		{
			name:  "empty-string-to-starlark",
			goVal: "",
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go("").Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				strVal, ok := starval.(starlark.String)
				if !ok {
					t.Fatalf("expected starlark.String, got %T", starval)
				}
				if string(strVal) != "" {
					t.Errorf("expected empty string, got %q", string(strVal))
				}
			},
		},
		{
			name:  "false-to-starlark",
			goVal: false,
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(false).Starlark(&starval); err != nil {
					t.Fatal(err)
				}
				boolVal, ok := starval.(starlark.Bool)
				if !ok {
					t.Fatalf("expected starlark.Bool, got %T", starval)
				}
				if bool(boolVal) != false {
					t.Errorf("expected false, got true")
				}
			},
		},
		{
			name:  "pointer-starlark",
			goVal: &struct{ Msg, Target string }{Msg: "hello", Target: "world"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlarkstruct.Struct
				if err := goToStarlark(goVal, &starval); err != nil {
					t.Fatal(err)
				}
				val, err := starval.Attr("Msg")
				if err != nil {
					t.Fatalf("failed to get value from starlarkstruct.Struct: [%#v], %s ", starval, err)
				}

				if val.String() != `"hello"` {
					t.Errorf("unexpected value for starlark.Dict value: %s", val.String())
				}

			},
		},
		{
			name:  "list-string",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected list length %d", starval.Len())
				}
				if starval.Index(1).String() != `"World!"` {
					t.Errorf("unexpected list value: %s", starval.Index(1).String())
				}
			},
		},
		{
			name:  "list-pointer",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected list length %d", starval.Len())
				}
				if starval.Index(1).String() != `"World!"` {
					t.Errorf("unexpected list value: %s", starval.Index(1).String())
				}
			},
		},
		{
			name:  "list-numeric",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 3 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}

				intVal, _ := starval.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "list-value",
			goVal: []int{1, 2, math.MaxInt8},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.Value
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				list := starval.(*starlark.List)
				if list.Len() != 3 {
					t.Errorf("unexpected tuple length %d", list.Len())
				}

				intVal, _ := list.Index(2).(starlark.Int).Int64()
				if intVal != math.MaxInt8 {
					t.Errorf("unexpected int value: %d", intVal)
				}
			},
		},
		{
			name:  "set-string",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Set
				if err := Go(goVal).StarlarkSet(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected set length %d", starval.Len())
				}
				iter := starval.Iterate()
				var val starlark.Value
				for iter.Next(&val) {
					if val.String() == `"World!"` {
						iter.Done()
						return
					}
				}
				t.Errorf("Stararlk set value not found")
			},
		},
		{
			name:  "set-pointer",
			goVal: []string{"Hello", "World!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval *starlark.Set
				if err := Go(goVal).StarlarkSet(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 2 {
					t.Errorf("unexpected set length %d", starval.Len())
				}
				iter := starval.Iterate()
				var val starlark.Value
				for iter.Next(&val) {
					if val.String() == `"World!"` {
						iter.Done()
						return
					}
				}
				t.Errorf("Stararlk set value not found")
			},
		},
		{
			name:  "set-mix",
			goVal: []interface{}{1, 2, 3, "Go!"},
			eval: func(t *testing.T, goVal interface{}) {
				var starval starlark.List
				if err := Go(goVal).StarlarkList(&starval); err != nil {
					t.Fatal(err)
				}
				if starval.Len() != 4 {
					t.Errorf("unexpected tuple length %d", starval.Len())
				}
				iter := starval.Iterate()
				var val starlark.Value
				for iter.Next(&val) {
					if val.String() == `"Go!"` {
						iter.Done()
						return
					}
				}
				t.Errorf("Stararlk set value not found")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.goVal)
		})
	}
}

func TestGo2Star2Go(t *testing.T) {
	data := struct {
		Name  string `name:"name"`
		Count int    `name:"count"`
	}{
		Name:  "test",
		Count: 24,
	}

	// Go -> Starlark
	var star starlarkstruct.Struct
	if err := Go(data).Starlark(&star); err != nil {
		t.Fatal(err)
	}
	nameVal, err := star.Attr("name")
	if err != nil {
		t.Fatal(err)
	}
	if nameVal.String() != `"test"` {
		t.Errorf("unexpected attr name value: %s", nameVal.String())
	}

	countVal, err := star.Attr("count")
	if err != nil {
		t.Fatal(err)
	}
	if countVal.(starlark.Int).BigInt().Int64() != 24 {
		t.Errorf("unexpected attr name value: %s", nameVal.String())
	}

	// Starlark -> Go
	goData := struct {
		Count int64  `name:"count"`
		Name  string `name:"name"`
	}{}

	if err := Starlark(&star).Go(&goData); err != nil {
		t.Errorf("conversion failed: %s", err)
	}

	if goData.Count != 24 {
		t.Errorf("unexpected go struct field value: %d", goData.Count)
	}
}

func TestToStarlarkValue(t *testing.T) {
	tests := []struct {
		name   string
		goVal  any
		check  func(t *testing.T, v starlark.Value)
		hasErr bool
	}{
		{
			name:  "nil",
			goVal: nil,
			check: func(t *testing.T, v starlark.Value) {
				if v != starlark.None {
					t.Errorf("expected None, got %v", v)
				}
			},
		},
		{
			name:  "bool-true",
			goVal: true,
			check: func(t *testing.T, v starlark.Value) {
				if v != starlark.True {
					t.Errorf("expected True, got %v", v)
				}
			},
		},
		{
			name:  "bool-false",
			goVal: false,
			check: func(t *testing.T, v starlark.Value) {
				if v != starlark.False {
					t.Errorf("expected False, got %v", v)
				}
			},
		},
		{
			name:  "int",
			goVal: 42,
			check: func(t *testing.T, v starlark.Value) {
				i, ok := v.(starlark.Int)
				if !ok {
					t.Fatalf("expected starlark.Int, got %T", v)
				}
				val, _ := i.Int64()
				if val != 42 {
					t.Errorf("expected 42, got %d", val)
				}
			},
		},
		{
			name:  "int64",
			goVal: int64(math.MaxInt64),
			check: func(t *testing.T, v starlark.Value) {
				i := v.(starlark.Int)
				val, _ := i.Int64()
				if val != math.MaxInt64 {
					t.Errorf("expected MaxInt64, got %d", val)
				}
			},
		},
		{
			name:  "uint64",
			goVal: uint64(math.MaxUint64),
			check: func(t *testing.T, v starlark.Value) {
				i := v.(starlark.Int)
				val, _ := i.Uint64()
				if val != math.MaxUint64 {
					t.Errorf("expected MaxUint64, got %d", val)
				}
			},
		},
		{
			name:  "float64-fractional",
			goVal: 3.14,
			check: func(t *testing.T, v starlark.Value) {
				f, ok := v.(starlark.Float)
				if !ok {
					t.Fatalf("expected starlark.Float, got %T", v)
				}
				if float64(f) != 3.14 {
					t.Errorf("expected 3.14, got %v", f)
				}
			},
		},
		{
			name:  "float64-integer-semantics",
			goVal: 3.0,
			check: func(t *testing.T, v starlark.Value) {
				// JSON number semantics: 3.0 should become starlark.Int
				i, ok := v.(starlark.Int)
				if !ok {
					t.Fatalf("expected starlark.Int for 3.0, got %T (%v)", v, v)
				}
				val, _ := i.Int64()
				if val != 3 {
					t.Errorf("expected 3, got %d", val)
				}
			},
		},
		{
			name:  "string",
			goVal: "hello",
			check: func(t *testing.T, v starlark.Value) {
				s, ok := v.(starlark.String)
				if !ok {
					t.Fatalf("expected starlark.String, got %T", v)
				}
				if string(s) != "hello" {
					t.Errorf("expected hello, got %s", s)
				}
			},
		},
		{
			name:  "slice-any",
			goVal: []any{1, "two", true, nil},
			check: func(t *testing.T, v starlark.Value) {
				list, ok := v.(*starlark.List)
				if !ok {
					t.Fatalf("expected *starlark.List, got %T", v)
				}
				if list.Len() != 4 {
					t.Fatalf("expected 4 elements, got %d", list.Len())
				}
				if list.Index(3) != starlark.None {
					t.Errorf("expected None at index 3, got %v", list.Index(3))
				}
			},
		},
		{
			name:  "map-string-any",
			goVal: map[string]any{"a": 1, "b": "two"},
			check: func(t *testing.T, v starlark.Value) {
				dict, ok := v.(*starlark.Dict)
				if !ok {
					t.Fatalf("expected *starlark.Dict, got %T", v)
				}
				if dict.Len() != 2 {
					t.Fatalf("expected 2 entries, got %d", dict.Len())
				}
			},
		},
		{
			name:  "nested-map-with-slices",
			goVal: map[string]any{"items": []any{1, 2, 3}, "name": "test"},
			check: func(t *testing.T, v starlark.Value) {
				dict := v.(*starlark.Dict)
				items, _, _ := dict.Get(starlark.String("items"))
				list := items.(*starlark.List)
				if list.Len() != 3 {
					t.Errorf("expected 3 items, got %d", list.Len())
				}
			},
		},
		{
			name:  "starlark-value-passthrough",
			goVal: starlark.String("already starlark"),
			check: func(t *testing.T, v starlark.Value) {
				s, ok := v.(starlark.String)
				if !ok {
					t.Fatalf("expected starlark.String, got %T", v)
				}
				if string(s) != "already starlark" {
					t.Errorf("expected 'already starlark', got %s", s)
				}
			},
		},
		{
			name:  "typed-slice-strings",
			goVal: []string{"a", "b", "c"},
			check: func(t *testing.T, v starlark.Value) {
				list, ok := v.(*starlark.List)
				if !ok {
					t.Fatalf("expected *starlark.List, got %T", v)
				}
				if list.Len() != 3 {
					t.Errorf("expected 3 elements, got %d", list.Len())
				}
			},
		},
		{
			name:  "typed-map-string-int",
			goVal: map[string]int{"a": 1, "b": 2},
			check: func(t *testing.T, v starlark.Value) {
				dict, ok := v.(*starlark.Dict)
				if !ok {
					t.Fatalf("expected *starlark.Dict, got %T", v)
				}
				if dict.Len() != 2 {
					t.Errorf("expected 2 entries, got %d", dict.Len())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := Go[any](tt.goVal).ToStarlarkValue()
			if (err != nil) != tt.hasErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil {
				tt.check(t, v)
			}
		})
	}
}

func TestToStarlarkValue_Unsupported(t *testing.T) {
	_, err := Go(struct{ X int }{42}).ToStarlarkValue()
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}
}

func TestGoValue_TypeSpecific(t *testing.T) {
	// ToBool
	b, err := Go(true).ToBool()
	if err != nil {
		t.Fatal(err)
	}
	if b != starlark.True {
		t.Errorf("expected True, got %v", b)
	}

	_, err = Go(42).ToBool()
	if err == nil {
		t.Fatal("expected error for ToBool on int")
	}

	// ToInt
	i, err := Go(42).ToInt()
	if err != nil {
		t.Fatal(err)
	}
	val, _ := i.Int64()
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}

	_, err = Go("hello").ToInt()
	if err == nil {
		t.Fatal("expected error for ToInt on string")
	}

	// ToFloat
	f, err := Go(3.14).ToFloat()
	if err != nil {
		t.Fatal(err)
	}
	if float64(f) != 3.14 {
		t.Errorf("expected 3.14, got %v", f)
	}

	_, err = Go("hello").ToFloat()
	if err == nil {
		t.Fatal("expected error for ToFloat on string")
	}

	// ToString
	s, err := Go("hello").ToString()
	if err != nil {
		t.Fatal(err)
	}
	if string(s) != "hello" {
		t.Errorf("expected hello, got %s", s)
	}

	_, err = Go(42).ToString()
	if err == nil {
		t.Fatal("expected error for ToString on int")
	}
}

func TestGoValue_ToDict(t *testing.T) {
	// Map[string]int with sorted keys
	dict, err := Map(map[string]int{"b": 2, "a": 1, "c": 3}).ToDict()
	if err != nil {
		t.Fatal(err)
	}
	if dict.Len() != 3 {
		t.Fatalf("expected 3 entries, got %d", dict.Len())
	}
	// Verify keys are present
	v, found, _ := dict.Get(starlark.String("a"))
	if !found {
		t.Fatal("key 'a' not found")
	}
	iv, _ := v.(starlark.Int).Int64()
	if iv != 1 {
		t.Errorf("expected 1, got %d", iv)
	}

	// Error case
	_, err = Go(42).ToDict()
	if err == nil {
		t.Fatal("expected error for ToDict on int")
	}
}

func TestGoValue_ToList(t *testing.T) {
	list, err := Slice([]string{"x", "y", "z"}).ToList()
	if err != nil {
		t.Fatal(err)
	}
	if list.Len() != 3 {
		t.Fatalf("expected 3 elements, got %d", list.Len())
	}
	if list.Index(0).(starlark.String) != "x" {
		t.Errorf("expected 'x', got %v", list.Index(0))
	}

	// Error case
	_, err = Go(42).ToList()
	if err == nil {
		t.Fatal("expected error for ToList on int")
	}
}

func TestGoValue_ValueReturnsConcreteType(t *testing.T) {
	// Verify Value() returns the concrete type
	gv := Go(42)
	x := gv.Value() // This should compile — proves Value() returns int, not any
	if x != 42 {
		t.Errorf("expected 42, got %d", x)
	}
}

func TestToStarlarkValue_RoundTrip(t *testing.T) {
	// Test round-trip: map[string]any → Dict → map[string]any
	original := map[string]any{
		"name":   "test",
		"count":  float64(42), // JSON semantics: 42.0 → Int
		"active": true,
		"items":  []any{"a", "b"},
		"nested": map[string]any{"key": "value"},
	}

	starVal, err := Go[any](original).ToStarlarkValue()
	if err != nil {
		t.Fatal(err)
	}

	dict := starVal.(*starlark.Dict)
	goVal, err := Starlark(dict).ToGoValue()
	if err != nil {
		t.Fatal(err)
	}

	result := goVal.(map[string]any)
	if result["name"] != "test" {
		t.Errorf("name: expected test, got %v", result["name"])
	}
	if result["count"] != int64(42) {
		t.Errorf("count: expected 42, got %v (%T)", result["count"], result["count"])
	}
	if result["active"] != true {
		t.Errorf("active: expected true, got %v", result["active"])
	}
	items := result["items"].([]any)
	if len(items) != 2 || items[0] != "a" {
		t.Errorf("items: unexpected %v", items)
	}
	nested := result["nested"].(map[string]any)
	if nested["key"] != "value" {
		t.Errorf("nested.key: expected value, got %v", nested["key"])
	}
}
