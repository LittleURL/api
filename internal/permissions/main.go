package permissions

var RoleNames = []string{Admin, Editor, Viewer, Nobody}

type Role interface {
	// domain
	DomainRead() bool
	DomainWrite() bool

	// users
	UsersRead() bool
	UsersWrite() bool

	// links
	LinksRead() bool
	LinksWrite() bool
}
