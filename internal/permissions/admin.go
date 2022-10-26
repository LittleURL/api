package permissions

const Admin = "admin"

type AdminRole struct {}

// domains
func (role *AdminRole) DomainRead() bool  { return true }
func (role *AdminRole) DomainWrite() bool { return true }
