package entity

import "github.com/programemer/gin-admin/global"

type JwtBlacklist struct {
	global.GVA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
