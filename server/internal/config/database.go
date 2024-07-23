package config

import (
	"errors"
	"fmt"
)

var mysqlFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=60s"

var dbFormatMap = map[string]string{
	"mysql": mysqlFormat,
}

type database struct {
	Driver          string `yaml:"driver" json:"driver,omitempty"`
	Username        string `yaml:"username" json:"username,omitempty"`
	Password        string `yaml:"password" json:"password,omitempty"`
	Host            string `yaml:"host" json:"host,omitempty"`
	Port            int    `yaml:"port" json:"port,omitempty"`
	Dbname          string `yaml:"dbname" json:"dbname,omitempty"`
	ConnMaxIdleTime int    `yaml:"connMaxIdleTime" json:"connMaxIdleTime,omitempty"`
	ConnMaxLifeTime int    `yaml:"connMaxLifeTime" json:"connMaxLifeTime,omitempty"`
	MaxIdleConns    int    `yaml:"maxIdleConns" json:"maxIdleConns,omitempty"`
	MaxOpenConns    int    `yaml:"maxOpenConns" json:"maxOpenConns,omitempty"`
}

func (_self database) GetDsn() (dsn string, err error) {
	// 2-8	【实战】引入Gorm-掌握go语言map语法
	if format, ok := dbFormatMap[_self.Driver]; ok {
		dsn = fmt.Sprintf(format, _self.Username,
			_self.Password, _self.Host, _self.Port, _self.Dbname)
		return
	}
	err = errors.New("请检查 database.driver配置,不支持: " + _self.Driver)
	return
}
