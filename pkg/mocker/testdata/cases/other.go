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
