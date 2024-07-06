package ssh_term

import (
	"encoding/json"
	"fmt"
	"go-protector/server/internal/config"
	"os"
	"path"
	"strconv"
	"time"
)

var headerFmt = `{"version": 2, "width": %v, "height": %v, "timestamp": %v, "env": {"SHELL": "/bin/zsh", "TERM": "%v"}}`

type Record struct {
	file      *os.File
	timestamp time.Time
}

func NewRecord(term *Terminal) (_self *Record, err error) {
	recordPath := config.GetConfig().Server.RecordPath
	targetPath := path.Join(recordPath, strconv.FormatUint(term.ssoSessionId, 10))
	// 判断 recordPath 是否存在
	if _, err = os.Stat(targetPath); os.IsNotExist(err) {
		if err = os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return
		}
	}
	err = nil
	filePath := path.Join(targetPath, "record.cast")
	var file *os.File
	if file, err = os.Create(filePath); err != nil {
		return
	}
	// {"version": 2, "width": 119, "height": 26, "timestamp": 1719850008, "env": {"SHELL": "/bin/zsh", "TERM": "xterm-256color"}}
	_self = new(Record)
	_self.file = file
	_self.timestamp = term.ConnectAt
	headerByte := []byte(fmt.Sprintf(headerFmt, term.w, term.h, _self.timestamp.Unix(), term.term))
	headerByte = append(headerByte, []byte{'\n'}...)
	if _, err = file.Write(headerByte); err != nil {
		return
	}

	return
}

func (_self *Record) write(data *string) (err error) {
	if data == nil || len(*data) <= 0 {
		return
	}
	// [4.343478, "o", "l"]

	// 计算差值
	//var sub = float64(time.Now().Sub(_self.timestamp).Microseconds()) / float64(time.Millisecond)
	var sub = float32(time.Now().Sub(_self.timestamp).Seconds())
	var row []any
	row = append(row, sub, "o")
	row = append(row, data)
	var bytes []byte
	if bytes, err = json.Marshal(row); err != nil {
		return
	}

	bytes = append(bytes, []byte{'\n'}...)
	_, err = _self.file.Write(bytes)

	return
}
