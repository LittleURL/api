package permissions

import "github.com/deltabyte/littleurl-api/internal/entities"

const Admin = "admin"

type AdminRole struct {
	DomainID entities.DomainID
}

// domains
func (role *AdminRole) DomainRead(id string) bool  { return id == role.DomainID }
func (role *AdminRole) DomainWrite(id string) bool { return id == role.DomainID }
