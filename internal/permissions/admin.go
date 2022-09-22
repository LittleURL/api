package permissions

import "gitlab.com/deltabyte_/littleurl/api/internal/entities"

const Admin = "admin"

type AdminRole struct {
	DomainID entities.DomainID
}

// domains
func (role *AdminRole) DomainRead(id string) bool  { return id == role.DomainID }
func (role *AdminRole) DomainWrite(id string) bool { return id == role.DomainID }
