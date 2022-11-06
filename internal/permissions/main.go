package permissions

type Role interface {
	DomainRead() bool
	DomainWrite() bool
}
