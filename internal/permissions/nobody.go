package permissions

import "gitlab.com/deltabyte_/littleurl/api/internal/entities"

const Nobody = "nobody"

type NobodyRole struct {
	DomainID entities.DomainID
}

// domains
func (role *NobodyRole) DomainRead(id string) bool  { return false }
func (role *NobodyRole) DomainWrite(id string) bool { return false }
