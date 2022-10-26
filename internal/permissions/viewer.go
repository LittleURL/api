package permissions

const Viewer = "viewer"

type ViewerRole struct {}

// domains
func (role *ViewerRole) DomainRead() bool  { return true }
func (role *ViewerRole) DomainWrite() bool { return false }
