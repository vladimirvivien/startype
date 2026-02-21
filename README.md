# Startype

Startype provides two-way conversion between Go types and [Starlark-Go](https://github.com/google/starlark-go) API types. It supports both **typed target** conversion (reflection-based, into a specific Go/Starlark type) and **dynamic dispatch** conversion (type-switched, for `any` data like JSON or Kubernetes Unstructured objects).

## Features

* Generic types: `GoValue[T any]` and `StarValue[T starlark.Value]` with type-safe `Value()` accessors
* **Typed target** conversion: `Go(val).Starlark(&target)` / `Starlark(val).Go(&target)` (reflection-based)
* **Dynamic dispatch** conversion: `Go(val).ToStarlarkValue()` / `Starlark(val).ToGoValue()` (type-switched)
* **Type-specific converters**: `ToBool`, `ToInt`, `ToFloat`, `ToString` (both directions)
* **Container converters**: `ToDict`, `ToList`, `ToMap`, `ToSlice` with convenience constructors `Map()`, `Slice()`, `Dict()`, `List()`
* Convert Go `slice`, `array`, `map`, and `struct` types to compatible Starlark types
* Convert Starlark `Dict`, `StringDict`, `List`, `Set`, and `Struct` to compatible Go types
* Map Starlark keyword args to Go struct values via `Kwargs()`
* Map both positional and keyword args via `Args()` (replacement for `starlark.UnpackArgs`)
* Struct tag support: `name`, `position`, `required`, `optional`

## API Overview

### Constructors

```go
startype.Go(val)           // *GoValue[T] — wraps any Go value
startype.Starlark(val)     // *StarValue[T] — wraps any Starlark value
startype.Map(m)            // *GoValue[map[K]V] — alias for Go(m)
startype.Slice(s)          // *GoValue[[]T] — alias for Go(s)
startype.Dict(d)           // *StarValue[*starlark.Dict] — alias for Starlark(d)
startype.List(l)           // *StarValue[*starlark.List] — alias for Starlark(l)
startype.Kwargs(kwargs)    // keyword args processor
startype.Args(args, kwargs)// positional + keyword args processor
```

### Go to Starlark

```go
// Typed target (reflection-based)
var star starlark.Int
startype.Go(42).Starlark(&star)

// Dynamic dispatch (type-switched)
val, err := startype.Go(myAny).ToStarlarkValue()

// Type-specific
b, err := startype.Go(true).ToBool()       // starlark.Bool
n, err := startype.Go(42).ToInt()           // starlark.Int
f, err := startype.Go(3.14).ToFloat()       // starlark.Float
s, err := startype.Go("hi").ToString()      // starlark.String

// Containers
dict, err := startype.Map(map[string]int{"a": 1}).ToDict()   // *starlark.Dict
list, err := startype.Slice([]string{"x", "y"}).ToList()      // *starlark.List
```

### Starlark to Go

```go
// Typed target (reflection-based)
var msg string
startype.Starlark(starStr).Go(&msg)

// Dynamic dispatch (type-switched)
goVal, err := startype.Starlark(starVal).ToGoValue()

// Type-specific
ok, err := startype.Starlark(starBool).ToBool()      // bool
num, err := startype.Starlark(starInt).ToInt64()       // int64
f, err := startype.Starlark(starFloat).ToFloat64()     // float64
s, err := startype.Starlark(starStr).ToString()        // string

// Containers
m, err := startype.Dict(starDict).ToMap()              // map[string]any
s, err := startype.List(starList).ToSlice()             // []any
```

### Generic Value Access

`Value()` returns the concrete type, not `any`:

```go
gv := startype.Go(42)
x := gv.Value()                  // x is int (not any)

sv := startype.Starlark(myDict)
d := sv.Value()                  // d is *starlark.Dict (not starlark.Value)
```

## Examples

### Go struct to Starlark struct

```go
data := struct {
    Name  string
    Count int
}{Name: "test", Count: 24}

var star starlarkstruct.Struct
if err := startype.Go(data).Starlark(&star); err != nil {
    log.Fatal(err)
}

nameVal, _ := star.Attr("name")
fmt.Println(nameVal) // "test"
```

### Starlark dict to Go map

```go
dict := starlark.NewDict(2)
_ = dict.SetKey(starlark.String("msg0"), starlark.String("Hello"))
_ = dict.SetKey(starlark.String("msg1"), starlark.String("World!"))

gomap := make(map[string]string)
if err := startype.Starlark(dict).Go(&gomap); err != nil {
    log.Fatal(err)
}
fmt.Println(gomap["msg0"]) // Hello
```

### Dynamic dispatch (any data)

Useful for JSON unmarshal results, Kubernetes Unstructured objects, etc.:

```go
// Go any → Starlark value
data := map[string]any{
    "name": "test",
    "tags": []any{"a", "b"},
    "count": 42,
}
val, err := startype.Go(data).ToStarlarkValue()
// val is a *starlark.Dict with nested List and Int

// Starlark value → Go any
goVal, err := startype.Starlark(val).ToGoValue()
// goVal is map[string]any with nested []any and int64
```

### Struct tags

```go
data := struct {
    Name  string `name:"msg0"`
    Count int    `name:"msg1"`
}{Name: "test", Count: 24}

var star starlarkstruct.Struct
startype.Go(data).Starlark(&star)
nameVal, _ := star.Attr("msg0") // uses tag name
```

### Keyword argument processing

```go
func myBuiltin(thread *starlark.Thread, fn *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
    var params struct {
        Path     string `name:"path" position:"0" required:"true"`
        Encoding string `name:"encoding" position:"1"`
        Force    bool   `name:"force"`
    }
    if err := startype.Args(args, kwargs).Go(&params); err != nil {
        return nil, err
    }
    // params.Path, params.Encoding, params.Force are populated
}
```

### Struct tags for Args

| Tag | Example | Description |
|-----|---------|-------------|
| `name` | `name:"path"` | Keyword argument name |
| `position` | `position:"0"` | Positional argument index (0-based) |
| `required` | `required:"true"` | Argument must be provided |
| `optional` | `optional:"true"` | Argument may be omitted (for `Kwargs()`) |

A field can have both `name` and `position` tags to accept either calling style. If both provide a value, the keyword argument wins.

## Dynamic Dispatch Type Mapping

### Go to Starlark (`ToStarlarkValue`)

| Go type | Starlark type |
|---------|---------------|
| `nil` | `None` |
| `bool` | `Bool` |
| `int`, `int64`, etc. | `Int` |
| `float64` (no fractional part) | `Int` |
| `float64` (fractional) | `Float` |
| `string` | `String` |
| `[]any` | `List` (recursive) |
| `map[string]any` | `Dict` (sorted keys, recursive) |

### Starlark to Go (`ToGoValue`)

| Starlark type | Go type |
|---------------|---------|
| `None` | `nil` |
| `Bool` | `bool` |
| `Int` | `int64` |
| `Float` | `float64` |
| `String` | `string` |
| `List`, `Tuple` | `[]any` (recursive) |
| `Dict` | `map[string]any` (recursive, string keys required) |

For additional examples, see the test files.
