package startype

import (
	"fmt"
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestStarArgs2Go(t *testing.T) {
	tests := []struct {
		name   string
		args   starlark.Tuple
		kwargs []starlark.Tuple
		eval   func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple)
	}{
		{
			name:   "positional only",
			args:   starlark.Tuple{starlark.String("/path/to/file")},
			kwargs: nil,
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path     string `name:"path" position:"0"`
					Encoding string `name:"encoding" position:"1"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Path != "/path/to/file" {
					t.Errorf("unexpected path: %s", val.Path)
				}
				if val.Encoding != "" {
					t.Errorf("expected empty encoding, got: %s", val.Encoding)
				}
			},
		},
		{
			name: "kwargs only",
			args: nil,
			kwargs: []starlark.Tuple{
				{starlark.String("path"), starlark.String("/path/to/file")},
				{starlark.String("encoding"), starlark.String("utf-8")},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path     string `name:"path" position:"0"`
					Encoding string `name:"encoding" position:"1"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Path != "/path/to/file" {
					t.Errorf("unexpected path: %s", val.Path)
				}
				if val.Encoding != "utf-8" {
					t.Errorf("unexpected encoding: %s", val.Encoding)
				}
			},
		},
		{
			name: "mixed positional and kwargs",
			args: starlark.Tuple{starlark.String("/path/to/file")},
			kwargs: []starlark.Tuple{
				{starlark.String("encoding"), starlark.String("utf-8")},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path     string `name:"path" position:"0"`
					Encoding string `name:"encoding" position:"1"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Path != "/path/to/file" {
					t.Errorf("unexpected path: %s", val.Path)
				}
				if val.Encoding != "utf-8" {
					t.Errorf("unexpected encoding: %s", val.Encoding)
				}
			},
		},
		{
			name: "kwarg overrides positional",
			args: starlark.Tuple{starlark.String("/positional")},
			kwargs: []starlark.Tuple{
				{starlark.String("path"), starlark.String("/keyword")},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path string `name:"path" position:"0"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Path != "/keyword" {
					t.Errorf("expected kwarg to override, got: %s", val.Path)
				}
			},
		},
		{
			name:   "missing required arg",
			args:   nil,
			kwargs: nil,
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path string `name:"path" position:"0" required:"true"`
				}
				if err := Args(args, kwargs).Go(&val); err == nil {
					t.Fatal("expected error for missing required arg")
				}
			},
		},
		{
			name: "unknown kwarg",
			args: nil,
			kwargs: []starlark.Tuple{
				{starlark.String("unknown"), starlark.String("value")},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path string `name:"path" position:"0"`
				}
				if err := Args(args, kwargs).Go(&val); err == nil {
					t.Fatal("expected error for unknown kwarg")
				}
			},
		},
		{
			name:   "extra positional arg",
			args:   starlark.Tuple{starlark.String("a"), starlark.String("b")},
			kwargs: nil,
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path string `name:"path" position:"0"`
				}
				if err := Args(args, kwargs).Go(&val); err == nil {
					t.Fatal("expected error for extra positional arg")
				}
			},
		},
		{
			name: "complex types - list and int",
			args: starlark.Tuple{
				starlark.NewList([]starlark.Value{
					starlark.String("alice"),
					starlark.String("bob"),
				}),
			},
			kwargs: []starlark.Tuple{
				{starlark.String("count"), starlark.MakeInt(42)},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Names []string `name:"names" position:"0"`
					Count int64    `name:"count" position:"1"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if len(val.Names) != 2 || val.Names[0] != "alice" {
					t.Errorf("unexpected names: %v", val.Names)
				}
				if val.Count != 42 {
					t.Errorf("unexpected count: %d", val.Count)
				}
			},
		},
		{
			name: "optional and required mixed",
			args: starlark.Tuple{starlark.String("/path")},
			kwargs: []starlark.Tuple{
				{starlark.String("count"), starlark.MakeInt(10)},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path     string `name:"path" position:"0" required:"true"`
					Encoding string `name:"encoding" position:"1"` // optional
					Count    int64  `name:"count" position:"2" required:"true"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Path != "/path" {
					t.Errorf("unexpected path: %s", val.Path)
				}
				if val.Encoding != "" {
					t.Errorf("expected empty encoding: %s", val.Encoding)
				}
				if val.Count != 10 {
					t.Errorf("unexpected count: %d", val.Count)
				}
			},
		},
		{
			name: "pointer field types",
			args: starlark.Tuple{starlark.String("value")},
			kwargs: []starlark.Tuple{
				{starlark.String("num"), starlark.MakeInt(100)},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Str *string `name:"str" position:"0"`
					Num *int64  `name:"num" position:"1"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Str == nil || *val.Str != "value" {
					t.Errorf("unexpected str: %v", val.Str)
				}
				if val.Num == nil || *val.Num != 100 {
					t.Errorf("unexpected num: %v", val.Num)
				}
			},
		},
		{
			name:   "empty args and kwargs",
			args:   nil,
			kwargs: nil,
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path string `name:"path" position:"0"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Path != "" {
					t.Errorf("expected empty path, got: %s", val.Path)
				}
			},
		},
		{
			name: "bool and float types",
			args: starlark.Tuple{starlark.Bool(true)},
			kwargs: []starlark.Tuple{
				{starlark.String("rate"), starlark.Float(3.14)},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Enabled bool    `name:"enabled" position:"0"`
					Rate    float64 `name:"rate" position:"1"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if !val.Enabled {
					t.Errorf("expected enabled to be true")
				}
				if val.Rate != 3.14 {
					t.Errorf("unexpected rate: %f", val.Rate)
				}
			},
		},
		{
			name: "dict type",
			args: nil,
			kwargs: []starlark.Tuple{
				{starlark.String("env"), func() *starlark.Dict {
					d := starlark.NewDict(2)
					_ = d.SetKey(starlark.String("HOME"), starlark.String("/home/user"))
					_ = d.SetKey(starlark.String("PATH"), starlark.String("/usr/bin"))
					return d
				}()},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Env map[string]string `name:"env" position:"0"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Env["HOME"] != "/home/user" {
					t.Errorf("unexpected HOME: %s", val.Env["HOME"])
				}
				if val.Env["PATH"] != "/usr/bin" {
					t.Errorf("unexpected PATH: %s", val.Env["PATH"])
				}
			},
		},
		{
			name: "required yes syntax",
			args: nil,
			kwargs: []starlark.Tuple{
				{starlark.String("path"), starlark.String("/some/path")},
			},
			eval: func(t *testing.T, args starlark.Tuple, kwargs []starlark.Tuple) {
				var val struct {
					Path string `name:"path" position:"0" required:"yes"`
				}
				if err := Args(args, kwargs).Go(&val); err != nil {
					t.Fatal(err)
				}
				if val.Path != "/some/path" {
					t.Errorf("unexpected path: %s", val.Path)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.eval(t, test.args, test.kwargs)
		})
	}
}

func TestGoToArgsToGo(t *testing.T) {
	type FileInfo struct {
		Name string
		Size int64
	}
	type readArgs struct {
		Path string   `name:"path" position:"0"`
		Info FileInfo `name:"info" position:"1"`
	}

	// args+kwargs -> Go struct
	var arg readArgs
	infoStruct := starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		starlark.StringDict{
			"name": starlark.String("test.txt"),
			"size": starlark.MakeInt(1024),
		},
	)
	err := Args(
		starlark.Tuple{starlark.String("/path/to/file")},
		[]starlark.Tuple{
			{starlark.String("info"), infoStruct},
		},
	).Go(&arg)
	if err != nil {
		t.Fatal(err)
	}

	// Go struct -> Starlark struct
	starOut := new(starlarkstruct.Struct)
	if err := Go(arg).Starlark(starOut); err != nil {
		t.Fatal(fmt.Errorf("Go->Starlark error: %v", err))
	}

	// Starlark struct -> Go struct
	var arg2 readArgs
	if err := Starlark(starOut).Go(&arg2); err != nil {
		t.Fatal(fmt.Errorf("Starlark->Go error: %w", err))
	}

	if arg2.Path != "/path/to/file" {
		t.Errorf("unexpected path: %s", arg2.Path)
	}
	if arg2.Info.Name != "test.txt" {
		t.Errorf("unexpected info.name: %s", arg2.Info.Name)
	}
	if arg2.Info.Size != 1024 {
		t.Errorf("unexpected info.size: %d", arg2.Info.Size)
	}
}

func TestGoToArgsToGo2(t *testing.T) {
	type execArgs struct {
		Command string            `name:"command" position:"0"`
		Env     map[string]string `name:"env" position:"1"`
		Timeout int64             `name:"timeout" position:"2"`
	}

	// Build Starlark values
	envDict := starlark.NewDict(2)
	_ = envDict.SetKey(starlark.String("HOME"), starlark.String("/home/user"))
	_ = envDict.SetKey(starlark.String("PATH"), starlark.String("/usr/bin"))

	// args+kwargs -> Go struct
	var arg execArgs
	err := Args(
		starlark.Tuple{starlark.String("ls -la")},
		[]starlark.Tuple{
			{starlark.String("env"), envDict},
			{starlark.String("timeout"), starlark.MakeInt(30)},
		},
	).Go(&arg)
	if err != nil {
		t.Fatal(err)
	}

	// Go struct -> Starlark struct
	starOut := new(starlarkstruct.Struct)
	if err := Go(arg).Starlark(starOut); err != nil {
		t.Fatal(fmt.Errorf("Go->Starlark error: %v", err))
	}

	// Starlark struct -> Go struct
	var arg2 execArgs
	if err := Starlark(starOut).Go(&arg2); err != nil {
		t.Fatal(fmt.Errorf("Starlark->Go error: %w", err))
	}

	if arg2.Command != "ls -la" {
		t.Errorf("unexpected command: %s", arg2.Command)
	}
	if arg2.Env["HOME"] != "/home/user" {
		t.Errorf("unexpected env HOME: %s", arg2.Env["HOME"])
	}
	if arg2.Timeout != 30 {
		t.Errorf("unexpected timeout: %d", arg2.Timeout)
	}
}

func TestArgsInvalidInputs(t *testing.T) {
	t.Run("non-pointer dest", func(t *testing.T) {
		var val struct {
			Path string `name:"path" position:"0"`
		}
		err := Args(nil, nil).Go(val)
		if err == nil {
			t.Fatal("expected error for non-pointer dest")
		}
	})

	t.Run("nil pointer dest", func(t *testing.T) {
		var val *struct {
			Path string `name:"path" position:"0"`
		}
		err := Args(nil, nil).Go(val)
		if err == nil {
			t.Fatal("expected error for nil pointer dest")
		}
	})

	t.Run("non-struct dest", func(t *testing.T) {
		var val string
		err := Args(nil, nil).Go(&val)
		if err == nil {
			t.Fatal("expected error for non-struct dest")
		}
	})

	t.Run("invalid kwarg name type", func(t *testing.T) {
		var val struct {
			Path string `name:"path" position:"0"`
		}
		// Create kwarg with non-string name
		kwargs := []starlark.Tuple{
			{starlark.MakeInt(123), starlark.String("value")},
		}
		err := Args(nil, kwargs).Go(&val)
		if err == nil {
			t.Fatal("expected error for non-string kwarg name")
		}
	})
}
