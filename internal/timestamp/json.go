package timestamp

import (
	"time"
)

func (ts Timestamp) MarshalJSON() ([]byte, error) {
	return time.Time(ts).MarshalJSON()
}

func (ts *Timestamp) UnmarshalJSON(data []byte) error {
	t := &time.Time{}

	if err := t.UnmarshalJSON(data); err != nil {
		return err
	}

	*ts = Timestamp(*t)
	return nil
}
