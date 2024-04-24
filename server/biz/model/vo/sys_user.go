package vo

import (
	"go-protector/server/biz/model/entity"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type UserPage struct {
	ID            uint64   `json:"id"`
	LoginName     string   `json:"loginName"`
	Username      string   `json:"username"`
	Email         string   `json:"email"`
	Sex           string   `json:"sex"`
	Status        int      `json:"status"`
	DeptId        int      `json:"deptId"`
	DeptName      string   `json:"deptName"`
	PostIds       string   `json:"-"`
	PostNames     string   `json:"-"`
	PostIdSlice   []uint64 `gorm:"-" json:"postIds"`
	PostNameSlice []string `gorm:"-" json:"postNames"`
	RoleIds       string   `json:"-"`
	RoleNames     string   `json:"-"`
	RoleIdSlice   []uint64 `gorm:"-" json:"roleIds"`
	RoleNameSlice []string `gorm:"-" json:"roleNames"`
	entity.ModelDelete
}

func (_self *UserPage) AfterFind(tx *gorm.DB) (err error) {
	var tempId uint64
	if len(_self.PostIds) > 0 {
		if split := strings.Split(_self.PostIds, ","); len(split) > 0 {
			for _, idStr := range split {
				if tempId, err = strconv.ParseUint(idStr, 10, 64); err != nil {
					return
				}
				_self.PostIdSlice = append(_self.PostIdSlice, tempId)
			}
		}
	}
	if len(_self.PostNames) > 0 {
		if split := strings.Split(_self.PostNames, ","); len(split) > 0 {
			for _, name := range split {
				_self.PostNameSlice = append(_self.PostNameSlice, name)
			}
		}
	}
	if len(_self.RoleIds) > 0 {
		if split := strings.Split(_self.RoleIds, ","); len(split) > 0 {
			for _, idStr := range split {
				if tempId, err = strconv.ParseUint(idStr, 10, 64); err != nil {
					return
				}
				_self.RoleIdSlice = append(_self.RoleIdSlice, tempId)
			}
		}
	}
	if len(_self.RoleNames) > 0 {
		if split := strings.Split(_self.RoleNames, ","); len(split) > 0 {
			for _, name := range split {
				_self.RoleNameSlice = append(_self.RoleNameSlice, name)
			}
		}
	}

	return
}
