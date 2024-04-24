package c_type

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type RelationType string
type LoginPolicyCode string

type Time sql.NullTime

func (_self *Time) UnmarshalJSON(data []byte) (err error) {
	val := strings.ReplaceAll(string(data), "\"", "")
	if len(val) <= 0 || val == "null" {
		_self.Valid = false
		return nil
	}
	var now time.Time
	var layout string
	if len(strings.Split(val, " ")) <= 1 {
		layout = time.DateOnly

	} else {
		layout = time.DateTime
	}
	if now, err = time.Parse(layout, val); err != nil {
		return
	}
	_self.Valid = true
	_self.Time = now
	return

}

func (_self Time) MarshalJSON() ([]byte, error) {
	if !_self.Valid {
		return []byte("null"), nil
	}
	sprintf := fmt.Sprintf("\"%s\"", _self.Time.Format(time.DateTime))
	return []byte(sprintf), nil
}

func (_self *Time) String() string {
	if !_self.Valid {
		return "null"
	}
	return _self.Time.Format(time.DateTime)
}

// Value 天坑 https://github.com/golang/go/blob/master/src/database/sql/driver/types.go#L242
func (_self Time) Value() (driver.Value, error) {
	if !_self.Valid {
		return nil, nil
	}
	return _self.Time, nil
}

func (_self *Time) Scan(v interface{}) error {
	return (*sql.NullTime)(_self).Scan(v)
}

func NewTime(v time.Time) Time {
	if v.IsZero() {
		return Time{Valid: false}
	}
	return Time{Valid: true, Time: v}
}

func NowTime() Time {
	return NewTime(time.Now())
}
