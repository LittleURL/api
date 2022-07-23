package entities

import (
	"encoding/json"
	"fmt"
	"strings"
)

type AppMetadata struct {
	// uses a pointer so that the value is nullable, because auth0 uses null to remove a key
	Domains map[DomainID]*string `json:"domains,omitempty"`
}

// from { domains: { foo: "bar" } }
// to   { domain_foo: "bar" }
func (appMeta AppMetadata) MarshalJSON() ([]byte, error) {
	attr := make(map[string]*string)

	// domains
	for k, v := range appMeta.Domains {
		attr[fmt.Sprintf("domain_%s", k)] = v
	}

	attrJson, err := json.Marshal(attr)
	if err != nil {
		return nil, err
	}

	return attrJson, nil
}

// from { domain_foo: "bar" }
// to   { domains: { foo: "bar" } }
func (appMeta *AppMetadata) UnmarshalJSON(b []byte) error {
	attr := make(map[string]*string)
	if err := json.Unmarshal(b, &attr); err != nil {
		return err
	}

	// domains
	domains := make(map[string]*string)
	for k, v := range attr {
		if strings.HasPrefix(k, "domain_") {
			cleanName := strings.TrimPrefix(k, "domain_")
			domains[cleanName] = v
		}
	}
	appMeta.Domains = domains

	return nil
}
