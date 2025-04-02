package entity

import "gin-adminn/global"

type SysApi struct {
	global.GVA_MODEL
	Path        string `json:"path" gorm:"comment:api路径"`
	Description string `json:"description" gorm:"comment:api中文描述"`
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`
	Method      string `json:"method" gorm:"default:POST;comment:方法"`
}

func (SysApi) TableName() string {
	return "sys_api"
}

type SysIgnoreApi struct {
	global.GVA_MODEL
	Path   string `json:"path" gorm:"comment:api路径"`
	Method string `json:"method" gorm:"default:POST;comment:方法"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Flag   bool   `json:"flag"  gorm:"-"`
}

func (SysIgnoreApi) TableName() string {
	return "sys_ignore_api"
}
