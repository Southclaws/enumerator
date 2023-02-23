# enumerator

> Generate safe and validated enumerated types.

This package generates **safe** enumerated values from simple constants. The
generated values are structs with one unexported string field which prevents you
from instantiating them with arbitrary values. This ensures that when you come
across one of these types in your codebase, you can rest assured it is valid.

More than a basic [stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer)
as it's designed to work with existing string constants instead of ints and iota
and it generates the code to validate arbitrary input and output a typed value.

## TL;DR

Simply add the go:generate line and suffix any enum string types with `Enum`:

```go
package mypackage

//go:generate go run -mod=mod github.com/Southclaws/enumerator

type statusEnum string

const (
    success   statusEnum = "success"
    failure   statusEnum = "failure"
    inBetween statusEnum = "inbetween"
    notSure   statusEnum = "notsure"
)
```

And you'll get a validator function generated for you!

## Why?

Have you ever used a string type to represent some finite set of states?

```go
type Status string

const (
    Success   Status = "success"
    Failure   Status = "failure"
    InBetween Status = "inbetween"
    NotSure   Status = "notsure"
)
```

This is useful for making sure you pass the right value to the right place, but
where this falls short is that it doesn't provide you with enough confidence
that a value of type `Status` is _only_ ever going to be one of these 4 strings.

For example, you're completely free to do this:

```go
sneaky := Status("invalid!")
```

Because `Status` is just a simple type based on `string` and you can simply cast
any string value to this type and there's nothing built into the type system to
prevent this kind of behaviour.

A common solution to this is to declare constrained types as a struct with a
single unexported field which API consumers cannot set and thus cannot construct
arbitrary instances of the type. This means the owning package has control.

```go
type Status struct{v string}

var (
    Success   Status{"success"}
    Failure   Status{"failure"}
    InBetween Status{"inbetween"}
    NotSure   Status{"notsure"}
)
```

As well as the above code, you also need some form of conversion functions to
parse arbitrary values and validate them against the set of valid states and
also convert the values back into strings for easier usage.

That's what this package is all about.

## Usage

This tool scans the target package (which, by default is the current directory)
and generates the necessary code for constrained enumerated values for any type
suffixed with `Enum`.

You can use Go generate, and you don't need to supply any arguments:

```go
//go:generate go run -mod=mod github.com/Southclaws/enumerator
```

It will just run on the current package and find any types suffixed with `Enum`.

For example:

```go
type statusEnum string

const (
    success   statusEnum = "success"
    failure   statusEnum = "failure"
    inBetween statusEnum = "inbetween"
    notSure   statusEnum = "notsure"
)
```

Which will generate identical set of symbols except they are exported and use
the struct trick mentioned above to constrain instances to the finite set.

It will also generate a simple stringer:

```go
func (r Status) String() string {
    return string(r.v)
}
```

As well as a constructor which performs validation:

```go
func NewStatus(in string) (Status, error) {
    switch in {
    case "success":
        return Success, nil
    case "failure":
        return Failure, nil
    case "inBetween":
        return Inbetween, nil
    case "notSure":
        return Notsure, nil
    default:
        return Status{}, fmt.Errorf("invalid value for type 'Status': '%s'", in)
    }
}
```

And satisfies the `MarshalText`, `UnmarshalText`, `Value` and `Scan` interfaces.

And that's it! Super simple right now. Issues and pull requests welcome!
