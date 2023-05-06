package basemodel

import (
	"fmt"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

type BaseTime time.Time

func (t *BaseTime) UnmarshalJSON(data []byte) (err error) {
	newTime, err := time.ParseInLocation("\""+timeFormat+"\"", string(data), time.Local)
	*t = BaseTime(newTime)
	return
}

func (t BaseTime) MarshalJSON() ([]byte, error) {
	timeStr := fmt.Sprintf("\"%s\"", time.Time(t).Format(timeFormat))
	return []byte(timeStr), nil
}

func (t BaseTime) String() string {
	return time.Time(t).Format(timeFormat)
}
