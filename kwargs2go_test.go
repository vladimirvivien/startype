package startype

import (
	"fmt"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestStarKwargs2Go(t *testing.T) {
	tests := []struct {
		name   string
		kwargs []starlark.Tuple
		eval   func(t *testing.T, kwargs []starlark.Tuple)
	}{
		{
			name:   "missing explicit required arg",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				var val struct {
					A string `name:"a" required:"true"`
				}
				if err := Kwargs(kwargs).Go(&val); err == nil {
					t.Fatal("expecting required argument error")
				}
				if val.A != "" {
					t.Error("expecting empty value")
				}
			},
		},
		{
			name:   "missing explicit required arg",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a" required:"false"`
				}{}
				if err := Kwargs(kwargs).Go(&val); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
			},
		},
		{
			name:   "missing default required arg",
			kwargs: []starlark.Tuple{},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a"`
				}{}
				if err := Kwargs(kwargs).Go(&val); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
			},
		},
		{
			name: "all optional implied",
			kwargs: []starlark.Tuple{
				{starlark.String("a"), starlark.String("hello")},
				{starlark.String("b"), starlark.MakeInt(32)},
				{starlark.String("c"), starlark.MakeInt(64)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a"`
					B int64  `name:"b"`
				}{}
				if err := Kwargs(kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.A != "hello" {
					t.Errorf("unexpected value: %s", val.A)
				}
				if val.B != 32 {
					t.Errorf("unexpected value: %d", val.B)
				}
			},
		},
		{
			name: "all optional explicit",
			kwargs: []starlark.Tuple{
				{starlark.String("a"), starlark.String("hello")},
				{starlark.String("b"), starlark.MakeInt(32)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a" required:"false"`
					B int64  `name:"b" required:"false"`
				}{}
				if err := Kwargs(kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.A != "hello" {
					t.Errorf("unexpected value: %s", val.A)
				}
				if val.B != 32 {
					t.Errorf("unexpected value: %d", val.B)
				}
			},
		},
		{
			name: "optional and required",
			kwargs: []starlark.Tuple{
				{starlark.String("a"), starlark.String("hello")},
				{starlark.String("b"), starlark.MakeInt(32)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a" required:"true"`
					B int64  `name:"b" required:"false"`
				}{}
				if err := Kwargs(kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.A != "hello" {
					t.Errorf("unexpected value: %s", val.A)
				}
				if val.B != 32 {
					t.Errorf("unexpected value: %d", val.B)
				}
			},
		},
		{
			name: "implicit optional and explicit required",
			kwargs: []starlark.Tuple{
				{starlark.String("b"), starlark.MakeInt(32)},
			},
			eval: func(t *testing.T, kwargs []starlark.Tuple) {
				val := struct {
					A string `name:"a"`
					B int64  `name:"b" required:"true"`
				}{}
				if err := Kwargs(kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.A != "" {
					t.Errorf("unexpected value for `a`: %s", val.A)
				}
				if val.B != 32 {
					t.Errorf("unexpected value for `b`: %d", val.B)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.kwargs)
		})
	}
}

func TestGoToKwargToGo(t *testing.T) {
	type Desc struct {
		Name string
	}
	type inarg struct {
		Name  Desc  `name:"name"`
		Count int64 `name:"count"`
	}

	// kwarg -> Go
	var arg inarg
	err := Kwargs(
		[]starlark.Tuple{
			{
				starlark.String("name"),
				starlarkstruct.FromStringDict(
					starlarkstruct.Default,
					starlark.StringDict{"name": starlark.String("Hello")},
				),
			},
			{starlark.String("count"), starlark.MakeInt(266)},
		}).Go(&arg)
	if err != nil {
		t.Fatal(err)
	}

	// Go -> starlark struct
	starOut := new(starlarkstruct.Struct)
	if err := Go(arg).Starlark(starOut); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %v", err))
	}

	// starlark struct -> Go
	var arg2 inarg
	if err := Starlark(starOut).Go(&arg2); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %w", err))
	}

	if arg2.Name.Name != "Hello" {
		t.Errorf("Unexpected value: %s", arg2.Name.Name)
	}
}

func TestGoToKwargToGo2(t *testing.T) {
	type Desc struct {
		Name string
	}
	type inarg struct {
		Name  Desc  `name:"name"`
		Count int64 `name:"count"`
		B     string
	}

	// kwarg -> Go
	var arg inarg
	desc := Desc{Name: "World"}
	descArg := new(starlarkstruct.Struct)
	if err := Go(desc).Starlark(descArg); err != nil {
		t.Fatal(err)
	}
	err := Kwargs(
		[]starlark.Tuple{
			{starlark.String("name"), descArg},
			{starlark.String("count"), starlark.MakeInt(266)},
		}).Go(&arg)
	if err != nil {
		t.Fatal(err)
	}

	// Go -> starlark struct
	starOut := new(starlarkstruct.Struct)
	if err := Go(arg).Starlark(starOut); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %v", err))
	}

	// starlark struct -> Go
	var arg2 inarg
	if err := Starlark(starOut).Go(&arg2); err != nil {
		t.Fatal(fmt.Errorf("conversion error: %w", err))
	}

	if arg2.Name.Name != "World" {
		t.Errorf("Unexpected value: %s", arg2.Name.Name)
	}
}
