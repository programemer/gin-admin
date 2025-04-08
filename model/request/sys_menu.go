package request

import (
	"github.com/programemer/gin-admin/global"
	"github.com/programemer/gin-admin/model/entity"
)

// AddMenuAuthorityInfo Add menu authority info structure
type AddMenuAuthorityInfo struct {
	Menus       []entity.SysBaseMenu `json:"menus"`
	AuthorityId uint                 `json:"authorityId"` // 角色ID
}

func DefaultMenu() []entity.SysBaseMenu {
	return []entity.SysBaseMenu{{
		GVA_MODEL: global.GVA_MODEL{ID: 1},
		ParentId:  0,
		Path:      "dashboard",
		Name:      "dashboard",
		Component: "view/dashboard/index.vue",
		Sort:      1,
		Meta: entity.Meta{
			Title: "仪表盘",
			Icon:  "setting",
		},
	}}
}
