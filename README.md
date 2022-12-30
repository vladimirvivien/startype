# Startype 🤩

If you have used the Starlark API to build tools, you know it can be awkward when going from Go types to Starlark-backed types and vice versa.
Startype makes it easy to automatically convert (roundtrip) between Go types and Starlark-Go API types.

## Features

* Two-way conversion between Go and Starlark-Go types
* Two-way conversion for primitive types (bool, numeric, string, etc)
* Conversion of Go container types (slice, array, map, struct)
* Conversion of Starlark container types (Dict, StringDict, List, Set)
* Support for type pointers and and empty interface (any) types

## Examples

### Convert Go value to Starlark value
The following converts a struct to a (comptible) `*starlarkstruct.Struct` value:

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

    // validate
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

Startype can easily convert a Starlark-API value into a standard Go value as shown below:

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

### Use struct annotation to target Starlark-Go type fields

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

	star := make.Struct
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

	if err := Starlark(&star).Go(&goData); err != nil {
		t.Errorf("conversion failed: %s", err)
	}

	if godata.Message != "World!" {
		log.Errorf("unexpected go struct field value: %s", godata.Message)
	}
}
```

## Starlark built-in function keywords processing
Startype makes it easy to process Starlark keyword arguments (passed as tuples in builtin functions) by automatically map the provided arguments to a Go struct value. For instance, 

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
	if err := KwargsToGo(kwargs, &val); err != nil {
		t.Fatal(err)
	}

    fmt.Printl(args.Message) // prints hello
}
```

An argument can be marked as optional to avoid error when it is not provided. For instance, 
if argument `cnt` is not provided in the `kwargs` tuple, function `KwargsToGo` will not reoprt an error.

```go

func main() {
    kwargs := []starlark.Tuple{
        {starlark.String("msg"), starlark.String("hello")},
    }
    var args struct {
        Message string   `name:"msg"`
        Count int64      `name:"cnt optional:true"`
    }
	if err := KwargsToGo(kwargs, &val); err != nil {
		t.Fatal(err)
	}

    fmt.Printl(args.Message) // prints hello
}
```

For additional conversion examples, see the test files in this module.