package startype

import "go.starlark.net/starlark"

// DictConvertible is implemented by custom Starlark types that can
// represent themselves as a *starlark.Dict for serialization.
// Types satisfying this interface are automatically handled by
// starlarkToGo and starlarkValueToGo, enabling generic operations
// like yaml.encode() to work with custom types.
type DictConvertible interface {
	starlark.Value
	ToDict() *starlark.Dict
}
