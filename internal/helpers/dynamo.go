package helpers

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

func ChangeTrackingExpression(builder expression.Builder, userId string) expression.Builder {
	// timestamp
	builder = builder.WithUpdate(expression.Set(
		expression.Name("updated_at"),
		expression.Name(NowISO3601()),
	))

	// user ID
	builder = builder.WithUpdate(expression.Set(
		expression.Name("updated_by"),
		expression.Name(userId),
	))

	return builder
}
