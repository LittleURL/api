package permissions

import (
	"strings"

	"github.com/deltabyte/littleurl-api/internal/entities"
)

type Role interface {
	DomainRead() bool
	DomainWrite() bool
}

func ParseScopes(claims map[string]string) *map[entities.DomainID]Role {
	roles := map[entities.DomainID]Role{}

	for key, value := range claims {
		// ignore non-domain claims
		if !strings.HasPrefix(key, "domain_") {
			continue
		}
		
		domain := strings.Replace(key, "_domain", "", 1)
		switch value {
		case Admin:
			roles[domain] = &AdminRole{}

		case Editor:
			roles[domain] = &EditorRole{}

		case Viewer:
			roles[domain] = &ViewerRole{}
		
		default:
			roles[domain] = &NobodyRole{}
		}
	}

	return &roles
}