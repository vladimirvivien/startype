# Startype 🤩

Startype makes it easy to automatically convert (two-way) between Go types and Starlark-Go API types.

## Features

* Two-way conversion between Go and Starlark-Go types
* Two-way conversion for primitive types like `bool`, `integer`, `float`, and `string` types
* Convert Go `slice`, `array`, `map`, and `struct` types to compatible Starlark types
* Convert Starlark `Dict`, `StringDict`, `List`, `Set`, and `StarlarkStruct` to compatible Go types
* Map Starlark keyword args (from built-in functions) to Go struct values
* Map both positional and keyword args to Go struct values (replacement for `starlark.UnpackArgs`)
* Support for type pointers and empty interface (any) types

## Examples

### Convert Go value to Starlark value
The following converts a Go `struct` to a (compatible) `*starlarkstruct.Struct` value:

```go
func main() {
    data := struct {
		Name  string
		Count int
	}{
		Name:  "test",
		Count: 24,
	}

	var star starlarkstruct.Struct
	if err := startype.Go(data).Starlark(&star); err != nil {
		log.Fatal(err)
	}

	nameVal, err := star.Attr("name")
	if err != nil {
		log.Fatal(err)
	}
	if nameVal.String() != `"test"` {
		log.Fatal("unexpected attr name value: %s", nameVal.String())
	}
}
```
### Convert Starlark value to Go value

Startype can easily convert a Starlark-API value into a standard Go value. The following 
converts a `starlark.Dict` dictionary value to a Go `map[string]string` value.

```go
func main() {
    dict := starlark.NewDict(2)
    dict.SetKey(starlark.String("msg0"), starlark.String("Hello"))
    dict.SetKey(starlark.String("msg1"), starlark.String("World!"))

    gomap := make(map[string]string)
    if err := startype.Starlark(val).Go(&gomap); err != nil {
        log.Fatalf("failed to convert starlark to go value: %s", err)
    }
	
    if gomap["msg0"] != "Hello" {
        log.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
    }
    if gomap["msg1"] != "World!" {
        log.Fatalf("unexpected map[msg] value: %v", gomap["msg"])
    }
}
```

### Use struct annotations to control conversion

Startype supports struct annotations to describe field names to target during conversion. For instance, the following example uses the provided struct tags when creating Starlark-Go values. 

```go
func main() {
    data := struct {
		Name  string `name:"msg0"`
		Count int    `name:"msg1"`
	}{
		Name:  "test",
		Count: 24,
	}

	var star starlarkstruct.Struct
	if err := startype.Go(data).Starlark(&star); err != nil {
		t.Fatal(err)
	}

    // starlark struct field created with annotated name
    nameVal, err := star.Attr("msg0")
	if err != nil {
		t.Fatal(err)
	}
	if nameVal.String() != `"test"` {
		t.Errorf("unexpected attr name value: %s", nameVal.String())
	}
}
```

Similarly, Startype can use struct tags to copy Starlark-API Go values during conversion to Go values, as shown below:

```go
func main() {
    dict := starlark.StringDict{
        "mymsg0": starlark.String("Hello"),
        "mymsg1": starlark.String("World!"),
    }
    star := starlarkstruct.FromStringDict(starlark.String("struct"), dict)
    
    var godata struct {
        Salutation string   `name:"mymsg0"`
        Message string      `name:"mymsg1"`
	}

	if err := startype.Starlark(&star).Go(&godata); err != nil {
		log.Fatalf("conversion failed: %s", err)
	}

	if godata.Message != "World!" {
		log.Fatalf("unexpected go struct field value: %s", godata.Message)
	}
}
```

## Starlark keyword argument processing
Startype makes it easy to capture and process Starlark keyword arguments (passed as tuples in [built-in functions](https://github.com/google/starlark-go/blob/master/doc/spec.md#built-in-functions)) by automatically map the provided arguments to a Go struct value. For instance, the following maps `kwargs` (a stand in for a actual keyword arguments) to Go struct `args`:

```go

func main() {
    kwargs := []starlark.Tuple{
        {starlark.String("msg"), starlark.String("hello")},
        {starlark.String("cnt"), starlark.MakeInt(32)},
    }
    var args struct {
        Message string   `name:"msg"`
        Count int64      `name:"cnt"`
    }
	if err := startype.Kwargs(kwargs).Go(&val); err != nil {
		t.Fatal(err)
	}

    fmt.Println(args.Message) // prints hello
}
```

An argument can be marked as optional to avoid error if it is not provided. For instance,
if argument `cnt` is not provided in the `kwargs` tuple, function `KwargsToGo` will not report an error.

```go

func main() {
    kwargs := []starlark.Tuple{
        {starlark.String("msg"), starlark.String("hello")},
    }
    var args struct {
        Message string   `name:"msg"`
        Count int64      `name:"cnt" optional:"true"`
    }
	if err := startype.Kwargs(kwargs).Go(&args); err != nil {
		log.Fatal(err)
	}

    fmt.Println(args.Message) // prints hello
}
```

## Combined positional and keyword argument processing

For more flexible argument handling (similar to `starlark.UnpackArgs`), Startype provides the `Args()` function that handles both positional arguments and keyword arguments. This is useful when implementing Starlark built-in functions that accept arguments in either style.

```go
func myBuiltin(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
    var params struct {
        Path     string `name:"path" position:"0" required:"true"`
        Encoding string `name:"encoding" position:"1"`
        Force    bool   `name:"force" position:"2"`
    }
    if err := startype.Args(args, kwargs).Go(&params); err != nil {
        return nil, err
    }

    // params.Path, params.Encoding, params.Force are now populated
    // ...
}
```

### Struct tags for Args

| Tag | Example | Description |
|-----|---------|-------------|
| `name` | `name:"path"` | Keyword argument name |
| `position` | `position:"0"` | Positional argument index (0-based) |
| `required` | `required:"true"` | Argument must be provided (positionally or by keyword) |

A field can have both `name` and `position` tags to accept either calling style. If both a positional argument and a keyword argument provide a value for the same field, the keyword argument wins.

### Examples

Positional arguments only:
```go
// Starlark call: my_func("/path/to/file", "utf-8")
var params struct {
    Path     string `name:"path" position:"0"`
    Encoding string `name:"encoding" position:"1"`
}
startype.Args(args, kwargs).Go(&params)
// params.Path = "/path/to/file", params.Encoding = "utf-8"
```

Keyword arguments only:
```go
// Starlark call: my_func(path="/path/to/file", encoding="utf-8")
var params struct {
    Path     string `name:"path" position:"0"`
    Encoding string `name:"encoding" position:"1"`
}
startype.Args(args, kwargs).Go(&params)
// params.Path = "/path/to/file", params.Encoding = "utf-8"
```

Mixed positional and keyword:
```go
// Starlark call: my_func("/path/to/file", encoding="utf-8")
var params struct {
    Path     string `name:"path" position:"0"`
    Encoding string `name:"encoding" position:"1"`
}
startype.Args(args, kwargs).Go(&params)
// params.Path = "/path/to/file", params.Encoding = "utf-8"
```

Required arguments:
```go
var params struct {
    Path     string `name:"path" position:"0" required:"true"`
    Encoding string `name:"encoding" position:"1"` // optional, defaults to zero value
}
if err := startype.Args(args, kwargs).Go(&params); err != nil {
    // Error if path not provided: "missing required argument: path"
}
```

For additional conversion examples, see the test functions in the test files.