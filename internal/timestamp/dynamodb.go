package timestamp

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	av "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (ts Timestamp) MarshalDynamoDBAttributeValue() (ddbTypes.AttributeValue, error) {
	return &ddbTypes.AttributeValueMemberN{
		Value: strconv.FormatInt(time.Time(ts).UnixMilli(), 10),
	}, nil
}

func (ts *Timestamp) UnmarshalDynamoDBAttributeValue(attr ddbTypes.AttributeValue) error {
	tv, ok := attr.(*ddbTypes.AttributeValueMemberN)
	if !ok {
		return &av.UnmarshalTypeError{
			Value: fmt.Sprintf("%T", attr),
			Type:  reflect.TypeOf((*Timestamp)(nil)),
		}
	}

	t, err := decodeDDBTimestamp(tv.Value)
	if err != nil {
		return err
	}

	*ts = Timestamp(t)
	return nil
}

func decodeDDBTimestamp(n string) (time.Time, error) {
	v, err := strconv.ParseInt(n, 10, 64)

	if err != nil {
		return time.Time{}, &av.UnmarshalError{
			Err: err, Value: n, Type: timeType,
		}
	}

	return time.UnixMilli(v), nil
}
