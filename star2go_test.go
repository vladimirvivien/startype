package startype

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestStarlarkToGo(t *testing.T) {
	tests := []struct {
		name    string
		starVal starlark.Value
		eval    func(*testing.T, starlark.Value)
	}{
		// Passthrough types - starlark.Value
		{
			name:    "value-passthrough-string",
			starVal: starlark.String("hello"),
			eval: func(t *testing.T, val starlark.Value) {
				var dest starlark.Value
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				if dest != starlark.String("hello") {
					t.Fatalf("unexpected value: %v", dest)
				}
				if dest.Type() != "string" {
					t.Fatalf("unexpected type: %s", dest.Type())
				}
			},
		},
		{
			name:    "value-passthrough-list",
			starVal: starlark.NewList([]starlark.Value{starlark.MakeInt(1), starlark.MakeInt(2)}),
			eval: func(t *testing.T, val starlark.Value) {
				var dest starlark.Value
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				if dest.Type() != "list" {
					t.Fatalf("unexpected type: %s", dest.Type())
				}
				list := dest.(*starlark.List)
				if list.Len() != 2 {
					t.Fatalf("unexpected list length: %d", list.Len())
				}
			},
		},
		// Passthrough types - starlark.Callable
		{
			name: "callable-passthrough",
			starVal: starlark.NewBuiltin("testfn", func(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
				return starlark.None, nil
			}),
			eval: func(t *testing.T, val starlark.Value) {
				var dest starlark.Callable
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				if dest.Name() != "testfn" {
					t.Fatalf("unexpected callable name: %s", dest.Name())
				}
			},
		},
		{
			name:    "callable-error-on-non-callable",
			starVal: starlark.String("not a function"),
			eval: func(t *testing.T, val starlark.Value) {
				var dest starlark.Callable
				err := Starlark(val).Go(&dest)
				if err == nil {
					t.Fatal("expected error for non-callable value")
				}
				if !strings.Contains(err.Error(), "not callable") {
					t.Fatalf("unexpected error message: %s", err)
				}
			},
		},
		// Passthrough types - starlark.Bytes
		{
			name:    "bytes-passthrough",
			starVal: starlark.Bytes("hello bytes"),
			eval: func(t *testing.T, val starlark.Value) {
				var dest starlark.Bytes
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				if dest != starlark.Bytes("hello bytes") {
					t.Fatalf("unexpected bytes value: %s", dest)
				}
			},
		},
		{
			name:    "bytes-to-slice",
			starVal: starlark.Bytes("hello bytes"),
			eval: func(t *testing.T, val starlark.Value) {
				var dest []byte
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				if string(dest) != "hello bytes" {
					t.Fatalf("unexpected []byte value: %s", dest)
				}
			},
		},
		{
			name:    "bytes-to-any",
			starVal: starlark.Bytes("hello bytes"),
			eval: func(t *testing.T, val starlark.Value) {
				var dest any
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				if string(dest.([]byte)) != "hello bytes" {
					t.Fatalf("unexpected any ([]byte) value: %v", dest)
				}
			},
		},
		{
			name:    "addressability",
			starVal: starlark.Bool(true),
			eval: func(t *testing.T, val starlark.Value) {
				var boolVar bool
				err := Starlark(val).Go(boolVar)
				if err == nil {
					t.Fatalf("expecting addressability error, got nil")
				}
			},
		},
		{
			name:    "bool",
			starVal: starlark.Bool(true),
			eval: func(t *testing.T, val starlark.Value) {
				var boolVar bool
				err := Starlark(val).Go(&boolVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if !boolVar {
					t.Fatalf("unexpected bool value: %t", boolVar)
				}
			},
		},
		{
			name:    "bool-ptr",
			starVal: starlark.Bool(true),
			eval: func(t *testing.T, val starlark.Value) {
				var boolVar *bool
				err := Starlark(val).Go(&boolVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if !(*boolVar) {
					t.Fatalf("unexpected *bool value: %t", *boolVar)
				}
			},
		},
		{
			name:    "bool-any",
			starVal: starlark.Bool(true),
			eval: func(t *testing.T, val starlark.Value) {
				var boolVar any
				err := Starlark(val).Go(&boolVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if !(boolVar.(bool)) {
					t.Fatalf("unexpected any (bool) value: %t", boolVar.(bool))
				}
			},
		},
		{
			name:    "int32",
			starVal: starlark.MakeInt(math.MaxInt32),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar int32
				err := Starlark(val).Go(&intVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxInt32 {
					t.Fatalf("unexpected int32 value: %d", intVar)
				}
			},
		},
		{
			name:    "int32-pointer",
			starVal: starlark.MakeInt(math.MaxInt32),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar *int32
				err := Starlark(val).Go(&intVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if *intVar != math.MaxInt32 {
					t.Fatalf("unexpected *int32 value: %d", *intVar)
				}
			},
		},
		{
			name:    "int-any",
			starVal: starlark.MakeInt(math.MaxInt32),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar any
				err := Starlark(val).Go(&intVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar.(int64) != math.MaxInt32 {
					t.Fatalf("unexpected any int32 value: %d", intVar.(int32))
				}
			},
		},
		{
			name:    "int64",
			starVal: starlark.MakeInt64(math.MaxInt64),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar int64
				err := Starlark(val).Go(&intVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxInt64 {
					t.Fatalf("unexpected int64 value: %d", intVar)
				}
			},
		},
		{
			name:    "uint64",
			starVal: starlark.MakeUint64(math.MaxUint64),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar uint64
				err := Starlark(val).Go(&intVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxUint64 {
					t.Fatalf("unexpected uint64 value: %d", intVar)
				}
			},
		},
		{
			name:    "float32",
			starVal: starlark.Float(math.MaxFloat32),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar float32
				err := Starlark(val).Go(&floatVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if floatVar != math.MaxFloat32 {
					t.Fatalf("unexpected float32 value: %f", floatVar)
				}
			},
		},
		{
			name:    "float32-pointer",
			starVal: starlark.Float(math.MaxFloat32),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar *float32
				err := Starlark(val).Go(&floatVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if *floatVar != math.MaxFloat32 {
					t.Fatalf("unexpected float32 value: %f", *floatVar)
				}
			},
		},
		{
			name:    "float64",
			starVal: starlark.Float(math.MaxFloat64),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar float64
				err := Starlark(val).Go(&floatVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if floatVar != math.MaxFloat64 {
					t.Fatalf("unexpected float64 value: %f", floatVar)
				}
			},
		},
		{
			name:    "float64-any",
			starVal: starlark.Float(math.MaxFloat64),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar any
				err := Starlark(val).Go(&floatVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if floatVar.(float64) != math.MaxFloat64 {
					t.Fatalf("unexpected float64 any value: %f", floatVar)
				}
			},
		},
		{
			name:    "string",
			starVal: starlark.String("Hello World!"),
			eval: func(t *testing.T, val starlark.Value) {
				var strVar string
				err := Starlark(val).Go(&strVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if strVar != "Hello World!" {
					t.Fatalf("unexpected string value: %s", strVar)
				}
			},
		},
		{
			name:    "string-pointer",
			starVal: starlark.String("Hello World!"),
			eval: func(t *testing.T, val starlark.Value) {
				var strVar *string
				err := Starlark(val).Go(&strVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if *strVar != "Hello World!" {
					t.Fatalf("unexpected string value: %s", *strVar)
				}
			},
		},
		{
			name:    "string-any",
			starVal: starlark.String("Hello World!"),
			eval: func(t *testing.T, val starlark.Value) {
				var strVar any
				err := Starlark(val).Go(&strVar)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if strVar.(string) != "Hello World!" {
					t.Fatalf("unexpected string value: %s", strVar.(string))
				}
			},
		},
		{
			name:    "list-string",
			starVal: starlark.NewList([]starlark.Value{starlark.String("Hello"), starlark.String("World!")}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]string, 0)
				err := Starlark(val).Go(&slice)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if strings.Join(slice, " ") != "Hello World!" {
					t.Fatalf("unexpected string value: %v", slice)
				}
			},
		},
		{
			name:    "list-numbers",
			starVal: starlark.NewList([]starlark.Value{starlark.MakeInt64(math.MaxInt64), starlark.MakeInt(math.MaxInt32)}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]int64, 0)
				err := Starlark(val).Go(&slice)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0] != math.MaxInt64 {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1] != math.MaxInt32 {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
			},
		},
		{
			name:    "list-mixed",
			starVal: starlark.NewList([]starlark.Value{starlark.String("HelloWorld!"), starlark.MakeInt(math.MaxInt32)}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := Starlark(val).Go(&slice)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name:    "tuple-mixed",
			starVal: starlark.Tuple([]starlark.Value{starlark.String("HelloWorld!"), starlark.MakeInt(math.MaxInt32)}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := Starlark(val).Go(&slice)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name: "dict[string]string",
			starVal: func() *starlark.Dict {
				dict := starlark.NewDict(2)
				if err := dict.SetKey(starlark.String("msg0"), starlark.String("Hello")); err != nil {
					panic(err)
				}
				if err := dict.SetKey(starlark.String("msg1"), starlark.String("World!")); err != nil {
					panic(err)
				}
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]string)
				err := Starlark(val).Go(&gomap)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gomap["msg0"] != "Hello" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
				if gomap["msg1"] != "World!" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
			},
		},
		{
			name: "dict[string]int",
			starVal: func() *starlark.Dict {
				dict := starlark.NewDict(2)
				if err := dict.SetKey(starlark.String("msg0"), starlark.MakeInt(math.MaxInt32)); err != nil {
					panic(err)
				}
				if err := dict.SetKey(starlark.String("msg1"), starlark.MakeInt64(math.MaxInt64)); err != nil {
					panic(err)
				}
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]int64)
				err := Starlark(val).Go(&gomap)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gomap["msg0"] != math.MaxInt32 {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
				if gomap["msg1"] != math.MaxInt64 {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
			},
		},
		{
			name: "dict[int]int",
			starVal: func() *starlark.Dict {
				dict := starlark.NewDict(2)
				if err := dict.SetKey(starlark.MakeInt(1), starlark.MakeInt(math.MaxInt32)); err != nil {
					panic(err)
				}
				if err := dict.SetKey(starlark.MakeInt(2), starlark.MakeInt64(math.MaxInt64)); err != nil {
					panic(err)
				}
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[int64]int64)
				if err := Starlark(val).Go(&gomap); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gomap[1] != math.MaxInt32 {
					t.Fatalf("unexpected map[msg] value: %v", gomap[1])
				}
				if gomap[2] != math.MaxInt64 {
					t.Fatalf("unexpected map[msg] value: %v", gomap[2])
				}
			},
		},
		{
			name: "dict[string]inner-map",
			starVal: func() *starlark.Dict {
				inner := starlark.NewDict(1)
				if err := inner.SetKey(starlark.String("type"), starlark.String("web")); err != nil {
					t.Fatal(err)
				}
				dict := starlark.NewDict(1)
				if err := dict.SetKey(starlark.String("labels"), inner); err != nil {
					panic(err)
				}
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]map[string]string)
				if err := Starlark(val).Go(&gomap); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}

				inner := gomap["labels"]
				if inner == nil {
					t.Fatal("inner map is nil")
				}

				if inner["type"] != "web" {
					t.Fatalf("unexpected value for inner map: %#v", inner)
				}
			},
		},
		{
			name: "dict[string]inner-struct",
			starVal: func() *starlark.Dict {
				inner := starlark.StringDict{
					"msg0": starlark.String("hello"),
					"msg1": starlark.String("world"),
				}
				dict := starlark.NewDict(1)
				if err := dict.SetKey(starlark.String("messages"), starlarkstruct.FromStringDict(starlark.String("struct"), inner)); err != nil {
					panic(err)
				}
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				type innerStruct struct {
					Msg0 string
					Msg1 string
				}
				gomap := make(map[string]innerStruct)
				if err := Starlark(val).Go(&gomap); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}

				inner := gomap["messages"]

				if !reflect.DeepEqual(inner, innerStruct{Msg0: "hello", Msg1: "world"}) {
					t.Fatalf("unexpected value for inner struct: %#v", inner)
				}
			},
		},
		{
			name: "dict[string]any",
			starVal: func() *starlark.Dict {
				inner := starlark.NewDict(1)
				if err := inner.SetKey(starlark.String("inner-type"), starlark.String("web")); err != nil {
					t.Fatal(err)
				}
				dict := starlark.NewDict(1)
				if err := dict.SetKey(starlark.String("inner"), inner); err != nil {
					panic(err)
				}
				if err := dict.SetKey(starlark.String("name"), starlark.String("app")); err != nil {
					panic(err)
				}
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]any)
				if err := Starlark(val).Go(&gomap); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}

				if gomap["name"] != "app" {
					t.Fatal("unexpected value for key 'name':", gomap["name"])
				}
				inner := gomap["inner"]
				if inner == nil {
					t.Fatal("inner map is nil")
				}
				innerMap := inner.(map[any]any)

				if innerMap["inner-type"] != "web" {
					t.Fatalf("unexpected value for inner map: %#v", inner)
				}
			},
		},
		{
			name: "set-string",
			starVal: func() *starlark.Set {
				set := starlark.NewSet(2)
				if err := set.Insert(starlark.String("HelloWorld!")); err != nil {
					panic(err)
				}
				if err := set.Insert(starlark.MakeInt(math.MaxInt32)); err != nil {
					panic(err)
				}
				return set
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := Starlark(val).Go(&slice)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name: "set-mixed",
			starVal: func() *starlark.Set {
				set := starlark.NewSet(2)
				if err := set.Insert(starlark.String("HelloWorld!")); err != nil {
					panic(err)
				}
				if err := set.Insert(starlark.MakeInt(math.MaxInt32)); err != nil {
					panic(err)
				}
				return set
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]interface{}, 0)
				err := Starlark(val).Go(&slice)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if slice[0].(string) != "HelloWorld!" {
					t.Fatalf("unexpected slice[0] value: %v", slice[0])
				}
				if slice[1].(int64) != math.MaxInt32 {
					t.Fatalf("unexpected slice[1] value: %v, want %d", slice[1], math.MaxInt32)
				}
			},
		},
		{
			name:    "none-to-any",
			starVal: starlark.None,
			eval: func(t *testing.T, val starlark.Value) {
				var dest any
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				if dest != nil {
					t.Fatalf("expected nil, got %v", dest)
				}
			},
		},
		{
			name:    "list-to-any",
			starVal: starlark.NewList([]starlark.Value{starlark.String("hello"), starlark.MakeInt(42)}),
			eval: func(t *testing.T, val starlark.Value) {
				var dest any
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				result, ok := dest.([]any)
				if !ok {
					t.Fatalf("expected []any, got %T", dest)
				}
				if len(result) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(result))
				}
				if result[0].(string) != "hello" {
					t.Fatalf("expected 'hello', got %v", result[0])
				}
				if result[1].(int64) != 42 {
					t.Fatalf("expected 42, got %v", result[1])
				}
			},
		},
		{
			name:    "tuple-to-any",
			starVal: starlark.Tuple([]starlark.Value{starlark.String("world"), starlark.MakeInt(99)}),
			eval: func(t *testing.T, val starlark.Value) {
				var dest any
				if err := Starlark(val).Go(&dest); err != nil {
					t.Fatalf("failed to convert: %s", err)
				}
				result, ok := dest.([]any)
				if !ok {
					t.Fatalf("expected []any, got %T", dest)
				}
				if len(result) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(result))
				}
				if result[0].(string) != "world" {
					t.Fatalf("expected 'world', got %v", result[0])
				}
				if result[1].(int64) != 99 {
					t.Fatalf("expected 99, got %v", result[1])
				}
			},
		},
		{
			name: "struct-strings",
			starVal: func() *starlarkstruct.Struct {
				dict := starlark.StringDict{
					"msg0": starlark.String("Hello"),
					"msg1": starlark.String("World!"),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				var gostruct struct{ Msg0, Msg1 string }
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Msg0 != "Hello" {
					t.Fatalf("unexpected struct value: %v", gostruct.Msg0)
				}
				if gostruct.Msg1 != "World!" {
					t.Fatalf("unexpected struct value: %v", gostruct.Msg1)
				}
			},
		},
		{
			name: "struct-numbers",
			starVal: func() *starlarkstruct.Struct {
				dict := starlark.StringDict{
					"smallInt": starlark.MakeInt(math.MaxInt32),
					"bigInt":   starlark.MakeInt64(math.MaxInt64),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				var gostruct struct{ SmallInt, BigInt int64 }
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.SmallInt != math.MaxInt32 {
					t.Fatalf("unexpected struct.SmallInt value: %v", gostruct.SmallInt)
				}
				if gostruct.BigInt != math.MaxInt64 {
					t.Fatalf("unexpected struct.BigInt value: %v", gostruct.BigInt)
				}
			},
		},
		{
			name: "struct-mixed",
			starVal: func() *starlarkstruct.Struct {
				dict := starlark.StringDict{
					"smallInt": starlark.MakeInt(math.MaxInt32),
					"msg1":     starlark.String("int 32 bits"),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				var gostruct struct {
					SmallInt int32
					Msg1     any
				}
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.SmallInt != math.MaxInt32 {
					t.Fatalf("unexpected struct.SmallInt value: %v", gostruct.SmallInt)
				}
				if gostruct.Msg1.(string) != "int 32 bits" {
					t.Fatalf("unexpected struct.Msg1 value: %v", gostruct.Msg1)
				}
			},
		},
		{
			name: "struct-inner-map",
			starVal: func() *starlarkstruct.Struct {
				val := starlark.NewDict(1)
				if err := val.SetKey(starlark.String("msg0"), starlark.String("Hello-World")); err != nil {
					panic(err)
				}
				if err := val.SetKey(starlark.String("msg1"), nil); err != nil {
					panic(err)
				}
				dict := starlark.StringDict{
					"key": starlark.String("test-message"),
					"val": val,
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				var gostruct struct {
					Key string
					Val map[string]interface{}
				}
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Key != "test-message" {
					t.Fatalf("unexpected value for struct.Key: %v", gostruct.Key)
				}
				if gostruct.Val == nil {
					t.Fatal("struct.Val is nil")
				}
				if gostruct.Val["msg0"].(string) != "Hello-World" {
					t.Fatalf("unexpected value from map")
				}
			},
		},
		{
			name: "struct-inner-struct",
			starVal: func() *starlarkstruct.Struct {
				inner := starlark.StringDict{
					"msg0": starlark.String("hello"),
					"msg1": starlark.String("world"),
				}
				dict := starlark.StringDict{
					"mystruct": starlarkstruct.FromStringDict(starlark.String("struct"), inner),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				type innerStruct struct {
					Msg0 string
					Msg1 string
				}
				var gostruct struct {
					Mystruct innerStruct
				}
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Mystruct.Msg0 != "hello" {
					t.Fatalf("unexpected value for struct.innerStruct.Msg0: %v", gostruct.Mystruct.Msg0)
				}
				if gostruct.Mystruct.Msg1 != "world" {
					t.Fatalf("unexpected value for struct.innerStruct.Msg1: %v", gostruct.Mystruct.Msg1)
				}
			},
		},
		{
			name: "struct-inner-structptr",
			starVal: func() *starlarkstruct.Struct {
				inner := starlark.StringDict{
					"msg0": starlark.String("hello"),
					"msg1": starlark.String("world"),
				}
				dict := starlark.StringDict{
					"mystruct": starlarkstruct.FromStringDict(starlark.String("struct"), inner),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				type InnerStruct struct {
					Msg0 string
					Msg1 string
				}
				var gostruct struct {
					Mystruct *InnerStruct
				}
				err := Starlark(val).Go(&gostruct)
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Mystruct.Msg0 != "hello" {
					t.Fatalf("unexpected value for struct.innerStruct.Msg0: %v", gostruct.Mystruct.Msg0)
				}
				if gostruct.Mystruct.Msg1 != "world" {
					t.Fatalf("unexpected value for struct.innerStruct.Msg1: %v", gostruct.Mystruct.Msg1)
				}
			},
		},
		{
			name: "struct-annotated",
			starVal: func() *starlarkstruct.Struct {
				dict := starlark.StringDict{
					"mymsg0": starlark.String("Hello"),
					"mymsg1": starlark.String("World!"),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				var gostruct struct {
					Msg0 string `name:"mymsg0"`
					Msg1 string `name:"mymsg1"`
				}
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Msg0 != "Hello" {
					t.Fatalf("unexpected annotated struct value: %v", gostruct.Msg0)
				}
				if gostruct.Msg1 != "World!" {
					t.Fatalf("unexpected annotated struct value: %v", gostruct.Msg1)
				}
			},
		},
		{
			name: "struct-inner-struct-annotated",
			starVal: func() *starlarkstruct.Struct {
				inner := starlark.StringDict{
					"msg0": starlark.String("hello"),
					"msg1": starlark.String("world"),
				}
				dict := starlark.StringDict{
					"simplestruct": starlarkstruct.FromStringDict(starlark.String("struct"), inner),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				type innerStruct struct {
					Msg0 string
					Msg1 string
				}
				var gostruct struct {
					Mystruct innerStruct `name:"simplestruct"`
				}
				err := starlarkToGo(val, reflect.ValueOf(&gostruct).Elem())
				if err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Mystruct.Msg0 != "hello" {
					t.Fatalf("unexpected value for struct.innerStruct.Msg0: %v", gostruct.Mystruct.Msg0)
				}
				if gostruct.Mystruct.Msg1 != "world" {
					t.Fatalf("unexpected value for struct.innerStruct.Msg1: %v", gostruct.Mystruct.Msg1)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.starVal)
		})
	}
}

func TestStarGo(t *testing.T) {
	tests := []struct {
		name    string
		starVal starlark.Value
		eval    func(*testing.T, starlark.Value)
	}{
		{
			name:    "bool",
			starVal: starlark.Bool(true),
			eval: func(t *testing.T, val starlark.Value) {
				var boolVar bool
				if err := Starlark(val).Go(&boolVar); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if !boolVar {
					t.Fatalf("unexpected bool value: %t", boolVar)
				}
			},
		},

		{
			name:    "int64",
			starVal: starlark.MakeInt64(math.MaxInt64),
			eval: func(t *testing.T, val starlark.Value) {
				var intVar int64
				if err := Starlark(val).Go(&intVar); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if intVar != math.MaxInt64 {
					t.Fatalf("unexpected int32 value: %d", intVar)
				}
			},
		},
		{
			name:    "float64",
			starVal: starlark.Float(math.MaxFloat64),
			eval: func(t *testing.T, val starlark.Value) {
				var floatVar float64
				if err := Starlark(val).Go(&floatVar); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if floatVar != math.MaxFloat64 {
					t.Fatalf("unexpected float64 value: %f", floatVar)
				}
			},
		},
		{
			name:    "list-string",
			starVal: starlark.NewList([]starlark.Value{starlark.String("Hello"), starlark.String("World!")}),
			eval: func(t *testing.T, val starlark.Value) {
				slice := make([]string, 0)
				if err := Starlark(val).Go(&slice); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if strings.Join(slice, " ") != "Hello World!" {
					t.Fatalf("unexpected string value: %v", slice)
				}
			},
		},
		{
			name: "dict[string]string",
			starVal: func() *starlark.Dict {
				dict := starlark.NewDict(2)
				if err := dict.SetKey(starlark.String("msg0"), starlark.String("Hello")); err != nil {
					panic(err)
				}
				if err := dict.SetKey(starlark.String("msg1"), starlark.String("World!")); err != nil {
					panic(err)
				}
				return dict
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				gomap := make(map[string]string)
				if err := Starlark(val).Go(&gomap); err != nil {
					t.Fatalf("failed to convert starlark to go value: %s", err)
				}
				if gomap["msg0"] != "Hello" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
				if gomap["msg1"] != "World!" {
					t.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
				}
			},
		},
		{
			name: "struct",
			starVal: func() *starlarkstruct.Struct {
				dict := starlark.StringDict{
					"msg0": starlark.String("Hello"),
					"msg1": starlark.String(""),
					"msg2": starlark.MakeInt64(12),
				}
				return starlarkstruct.FromStringDict(starlark.String("struct"), dict)
			}(),
			eval: func(t *testing.T, val starlark.Value) {
				var gostruct struct {
					Msg0, Msg1 string
					Msg2       int64
				}
				if err := Starlark(val).Go(&gostruct); err != nil {
					t.Errorf("failed to convert starlark to go value: %s", err)
				}
				if gostruct.Msg0 != "Hello" {
					t.Errorf("unexpected struct.Msg0 value: %v", gostruct.Msg0)
				}
				if gostruct.Msg1 != "" {
					t.Errorf("unexpected struct.Msg1 value: %v", gostruct.Msg1)
				}
				if gostruct.Msg2 != 12 {
					t.Errorf("unexpected struct.Msg2 value %v", gostruct.Msg2)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.starVal)
		})
	}
}

func TestToGoValue(t *testing.T) {
	tests := []struct {
		name    string
		starVal starlark.Value
		check   func(t *testing.T, v any)
		hasErr  bool
	}{
		{
			name:    "none",
			starVal: starlark.None,
			check:   func(t *testing.T, v any) {
				if v != nil {
					t.Errorf("expected nil, got %v", v)
				}
			},
		},
		{
			name:    "bool",
			starVal: starlark.True,
			check:   func(t *testing.T, v any) {
				if v != true {
					t.Errorf("expected true, got %v", v)
				}
			},
		},
		{
			name:    "int",
			starVal: starlark.MakeInt(42),
			check:   func(t *testing.T, v any) {
				if v.(int64) != 42 {
					t.Errorf("expected 42, got %v", v)
				}
			},
		},
		{
			name:    "float",
			starVal: starlark.Float(3.14),
			check:   func(t *testing.T, v any) {
				if v.(float64) != 3.14 {
					t.Errorf("expected 3.14, got %v", v)
				}
			},
		},
		{
			name:    "string",
			starVal: starlark.String("hello"),
			check:   func(t *testing.T, v any) {
				if v.(string) != "hello" {
					t.Errorf("expected hello, got %v", v)
				}
			},
		},
		{
			name:    "list",
			starVal: starlark.NewList([]starlark.Value{starlark.MakeInt(1), starlark.String("two")}),
			check:   func(t *testing.T, v any) {
				list := v.([]any)
				if len(list) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(list))
				}
				if list[0].(int64) != 1 {
					t.Errorf("expected 1, got %v", list[0])
				}
				if list[1].(string) != "two" {
					t.Errorf("expected two, got %v", list[1])
				}
			},
		},
		{
			name:    "tuple",
			starVal: starlark.Tuple{starlark.MakeInt(10), starlark.String("ten")},
			check:   func(t *testing.T, v any) {
				list := v.([]any)
				if len(list) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(list))
				}
				if list[0].(int64) != 10 {
					t.Errorf("expected 10, got %v", list[0])
				}
			},
		},
		{
			name: "dict",
			starVal: func() *starlark.Dict {
				d := starlark.NewDict(2)
				_ = d.SetKey(starlark.String("a"), starlark.MakeInt(1))
				_ = d.SetKey(starlark.String("b"), starlark.String("two"))
				return d
			}(),
			check: func(t *testing.T, v any) {
				m := v.(map[string]any)
				if m["a"].(int64) != 1 {
					t.Errorf("expected 1, got %v", m["a"])
				}
				if m["b"].(string) != "two" {
					t.Errorf("expected two, got %v", m["b"])
				}
			},
		},
		{
			name: "dict-non-string-key",
			starVal: func() *starlark.Dict {
				d := starlark.NewDict(1)
				_ = d.SetKey(starlark.MakeInt(1), starlark.String("one"))
				return d
			}(),
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := Starlark(tt.starVal).ToGoValue()
			if (err != nil) != tt.hasErr {
				t.Fatalf("unexpected error: %v", err)
			}
			if err == nil {
				tt.check(t, v)
			}
		})
	}
}

func TestStarValue_TypeSpecific(t *testing.T) {
	b, err := Starlark(starlark.True).ToBool()
	if err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("expected true")
	}

	_, err = Starlark(starlark.MakeInt(1)).ToBool()
	if err == nil {
		t.Fatal("expected error for ToBool on int")
	}

	i, err := Starlark(starlark.MakeInt(42)).ToInt64()
	if err != nil {
		t.Fatal(err)
	}
	if i != 42 {
		t.Errorf("expected 42, got %d", i)
	}

	_, err = Starlark(starlark.String("hello")).ToInt64()
	if err == nil {
		t.Fatal("expected error for ToInt64 on string")
	}

	f, err := Starlark(starlark.Float(2.71)).ToFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f != 2.71 {
		t.Errorf("expected 2.71, got %f", f)
	}

	_, err = Starlark(starlark.String("hello")).ToFloat64()
	if err == nil {
		t.Fatal("expected error for ToFloat64 on string")
	}

	s, err := Starlark(starlark.String("world")).ToString()
	if err != nil {
		t.Fatal(err)
	}
	if s != "world" {
		t.Errorf("expected world, got %s", s)
	}

	_, err = Starlark(starlark.MakeInt(1)).ToString()
	if err == nil {
		t.Fatal("expected error for ToString on int")
	}
}

func TestStarValue_ToMap(t *testing.T) {
	d := starlark.NewDict(2)
	_ = d.SetKey(starlark.String("name"), starlark.String("test"))
	_ = d.SetKey(starlark.String("count"), starlark.MakeInt(42))

	m, err := Dict(d).ToMap()
	if err != nil {
		t.Fatal(err)
	}
	if m["name"] != "test" {
		t.Errorf("expected test, got %v", m["name"])
	}
	if m["count"].(int64) != 42 {
		t.Errorf("expected 42, got %v", m["count"])
	}

	_, err = Starlark(starlark.String("hello")).ToMap()
	if err == nil {
		t.Fatal("expected error for ToMap on string")
	}
}

func TestStarValue_ToSlice(t *testing.T) {
	l := starlark.NewList([]starlark.Value{starlark.String("a"), starlark.MakeInt(1)})

	s, err := List(l).ToSlice()
	if err != nil {
		t.Fatal(err)
	}
	if len(s) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(s))
	}
	if s[0].(string) != "a" {
		t.Errorf("expected a, got %v", s[0])
	}

	_, err = Starlark(starlark.String("hello")).ToSlice()
	if err == nil {
		t.Fatal("expected error for ToSlice on string")
	}
}

func TestStarValue_ValueReturnsConcreteType(t *testing.T) {
	sv := Starlark(starlark.True)
	b := sv.Value()
	if b != starlark.True {
		t.Errorf("expected True, got %v", b)
	}

	d := starlark.NewDict(0)
	dv := Starlark(d)
	_ = dv.Value()
}

func TestStarValue_BackwardCompat(t *testing.T) {
	var msg string
	if err := Starlark(starlark.String("hello")).Go(&msg); err != nil {
		t.Fatal(err)
	}
	if msg != "hello" {
		t.Errorf("expected hello, got %s", msg)
	}

	var num int64
	if err := Starlark(starlark.MakeInt(42)).Go(&num); err != nil {
		t.Fatal(err)
	}
	if num != 42 {
		t.Errorf("expected 42, got %d", num)
	}
}
