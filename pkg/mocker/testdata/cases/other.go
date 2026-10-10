package cases

// Concrete represents a non-interface type.
type Concrete struct{}

// Alias represents struct type alias.
type Alias = Concrete

// Empty represents an interface with no methods.
type Empty interface{}

// EmptyAny represents an interface with no methods.
type EmptyAny any

// Other is used as a type in one of the cases.
type Other struct{}

// ParamOne is a parametrized type with a single parameter.
type ParamOne[A any] struct{ DataA A }

func (p *ParamOne[A]) Get() any { return p.DataA }

// ParamTwo is a parametrized type with two parameters.
type ParamTwo[A any, B any] struct {
	DataA A
	DataB B
}

// Generic represents an interface with a type parameter.
type Generic[T any] interface{ Get() T }

// EmbedGeneric represents an interface embedding an instantiated generic one.
type EmbedGeneric interface{ Generic[int] }
