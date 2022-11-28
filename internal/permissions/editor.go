package permissions

const Editor = "editor"

type EditorRole struct {}

// domains
func (role *EditorRole) DomainRead() bool  { return true }
func (role *EditorRole) DomainWrite() bool { return false }

// users
func (role *EditorRole) UsersRead() bool  { return true }
func (role *EditorRole) UsersWrite() bool { return false }

// links
func (role *EditorRole) LinksRead() bool  { return true }
func (role *EditorRole) LinksWrite() bool { return true }