package permissions

const Nobody = "nobody"

type NobodyRole struct {}

// domains
func (role *NobodyRole) DomainRead() bool  { return false }
func (role *NobodyRole) DomainWrite() bool { return false }

// users
func (role *NobodyRole) UsersRead() bool  { return false }
func (role *NobodyRole) UsersWrite() bool { return false }
