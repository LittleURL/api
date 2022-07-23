package permissions

import "github.com/deltabyte/littleurl-api/internal/entities"

const Viewer = "viewer"

type ViewerRole struct {
	DomainID entities.DomainID
}

// domains
func (role *ViewerRole) DomainRead(id string) bool  { return id == role.DomainID }
func (role *ViewerRole) DomainWrite(id string) bool { return false }
