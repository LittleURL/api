package permissions

import (
	"strings"

	"gitlab.com/deltabyte_/littleurl/api/internal/entities"
)

type Role interface {
	DomainRead(entities.DomainID) bool
	DomainWrite(entities.DomainID) bool
}

func ParseClaims(claims map[string]string) map[entities.DomainID]Role {
	roles := map[entities.DomainID]Role{}

	for key, value := range claims {
		// ignore non-domain claims
		if !strings.HasPrefix(key, "domain_") {
			continue
		}

		domain := strings.Replace(key, "_domain", "", 1)
		switch value {
		case Admin:
			roles[domain] = &AdminRole{DomainID: domain}

		case Editor:
			roles[domain] = &EditorRole{DomainID: domain}

		case Viewer:
			roles[domain] = &ViewerRole{DomainID: domain}

		default:
			roles[domain] = &NobodyRole{DomainID: domain}
		}
	}

	return roles
}
