package helpers

import(
	cognitoTypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func FlattenCognitoUserAttributes(attributes []cognitoTypes.AttributeType) map[string]string {
	flattened := map[string]string{}

	for _, attribute := range attributes {
		flattened[*attribute.Name] = *attribute.Value
	}

	return flattened
}