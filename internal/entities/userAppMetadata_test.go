package entities

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalJSON(t *testing.T) {
	assert := assert.New(t)

	domainValue := "bazqux"
	appMeta := AppMetadata{
		Domains: map[string]*string{
			"foobar":        &domainValue,
			"removedDomain": nil,
		},
	}

	appMetaJson, err := json.Marshal(appMeta)
	if err != nil {
		t.Errorf("Failed to marshal to JSON: %s", err)
	}

	assert.JSONEq(`{"domain_foobar":"bazqux","domain_removedDomain":null}`, string(appMetaJson))
}

func TestUnmarshalJSON(t *testing.T) {
	assert := assert.New(t)

	inputJson := `{"domain_foobar":"bazqux","domain_removedDomain":null}`

	appMeta := &AppMetadata{}
	if err := json.Unmarshal([]byte(inputJson), appMeta); err != nil {
		t.Errorf("Failed to unmarshal from JSON: %s", err)
	}

	domainValue := "bazqux"
	exampleAppMeta := &AppMetadata{
		Domains: map[string]*string{
			"foobar":        &domainValue,
			"removedDomain": nil,
		},
	}

	assert.Equal(exampleAppMeta, appMeta)
}
