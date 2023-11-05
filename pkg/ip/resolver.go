package ip

type Resolver interface {
	Resolve() (string, error)
}
