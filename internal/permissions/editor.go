package permissions

const Editor = "editor"

type EditorRole struct{}

// domains
func (role *EditorRole) DomainRead() bool  { return true }
func (role *EditorRole) DomainWrite() bool { return true }
