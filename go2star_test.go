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
				tuple := starval.(starlark.Tuple)
				if tuple.Len() != 3 {
					t.Errorf("unexpected tuple length %d", tuple.Len())
				}

				intVal, _ := tuple.Index(2).(starlark.Int).Int64()
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
				var starval starlark.Dict
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
				var starval starlark.Dict
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
				var starval starlark.Dict
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
				var starval starlark.Set
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
