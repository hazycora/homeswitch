package marshaltime

import (
	"encoding/json"
	time "time"
)

const (
	ISO8601Datetime = "2006-01-02T15:04:05.999Z"
	ISO8601Date     = "2006-01-02"
)

type Time time.Time
type Date time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	str := time.Time(t).Format(ISO8601Datetime)
	return json.Marshal(str)
}

func (t Date) MarshalJSON() ([]byte, error) {
	str := time.Time(t).Format(ISO8601Date)
	return json.Marshal(str)
}

func Now() Time {
	return Time(time.Now())
}
