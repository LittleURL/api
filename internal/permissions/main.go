package permissions

var RoleNames = []string{Admin, Editor, Viewer, Nobody}

type Role interface {
	DomainRead() bool
	DomainWrite() bool
	UsersRead() bool
	UsersWrite() bool
}
