package entities

import (
	"testing"

	"github.com/auth0/go-auth0/management"
	"github.com/stretchr/testify/assert"
)

func TestToAuth0User(t *testing.T) {
	assert := assert.New(t)

	domainString := "foobar"
	user := User{
		Id: "auth0|user-id-123",
		AppMetadata: AppMetadata{
			Domains: map[string]*string{
				"domain-id-123": &domainString,
			},
		},
	}

	auth0User, err := user.ToAuth0User()
	if err != nil {
		t.Error("Failed to convert user to auth0User")
	}

	exampleAuth0UserID := "auth0|user-id-123"
	exampleAuth0User := &management.User{
		ID: &exampleAuth0UserID,
		AppMetadata: &map[string]interface{}{
			"domain_domain-id-123": "foobar",
		},
	}

	assert.Equal(exampleAuth0User, auth0User)
}

func TestAuth0UserIDPrefix(t *testing.T) {
	// prefix is added
	user := User{Id: "user-id-123"}
	auth0User, _ := user.ToAuth0User()
	if auth0User.GetID() != "auth0|user-id-123" {
		t.Error("auth0UserID not prefixed")
	}

	// prefix already exists
	user2 := User{Id: "auth0|user-id-123"}
	auth0User2, _ := user2.ToAuth0User()
	if auth0User2.GetID() == "auth0|auth0|user-id-123" {
		t.Error("auth0UserID prefixed twice")
	}
}
