package permissions

const Nobody = "nobody"

type NobodyRole struct{}

// domains
func (role *NobodyRole) DomainRead() bool  { return false }
func (role *NobodyRole) DomainWrite() bool { return false }
