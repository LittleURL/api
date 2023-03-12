package timestamp

import (
	"reflect"
	"time"
)

type Timestamp time.Time

var timeType = reflect.TypeOf(time.Time{})

func Now() Timestamp {
	return Timestamp(time.Now())
}

func (ts Timestamp) String() string {
	return time.Time(ts).String()
}

func (ts *Timestamp) Touch() {
	*ts = Timestamp(time.Now())
}

func (ts *Timestamp) Until() time.Duration {
	return time.Until(time.Time(*ts))
}
