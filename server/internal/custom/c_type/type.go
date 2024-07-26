package c_type

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

var approveStatusSlice = []string{"待处理", "通过", "拒绝", "撤回", "超时未处理"}

type ApproveStatus int64

type ApproveType string

type SliceCondition string

type RelationType string
type LoginPolicyCode string

type MsgType int

type SessionStatus string

// 4-22	【实战】修改用户状态-掌握结构体更新与map更新的区别、自定义类型实现JSON序列化与数据库保存
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

func (_self ApproveStatus) String() string {
	var text string
	if _self > 4 || _self < 0 {
		return ""
	}
	text = approveStatusSlice[_self]
	return text
}
