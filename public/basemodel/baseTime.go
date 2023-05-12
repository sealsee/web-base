package basemodel

import (
	"database/sql/driver"
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

// gorm自定义数据类型须实现Scanner/Valuer接口
func (t BaseTime) Value() (driver.Value, error) {
	return time.Time(t).Format(timeFormat), nil
}

func (t *BaseTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = BaseTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
