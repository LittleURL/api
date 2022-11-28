package permissions

const Viewer = "viewer"

type ViewerRole struct{}

// domains
func (role *ViewerRole) DomainRead() bool  { return true }
func (role *ViewerRole) DomainWrite() bool { return false }

// users
func (role *ViewerRole) UsersRead() bool  { return false }
func (role *ViewerRole) UsersWrite() bool { return false }

// links
func (role *ViewerRole) LinksRead() bool  { return true }
func (role *ViewerRole) LinksWrite() bool { return false }
