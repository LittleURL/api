package permissions

const Admin = "admin"

type AdminRole struct {}

// domains
func (role *AdminRole) DomainRead() bool  { return true }
func (role *AdminRole) DomainWrite() bool { return true }

// users
func (role *AdminRole) UsersRead() bool  { return true }
func (role *AdminRole) UsersWrite() bool { return true }
