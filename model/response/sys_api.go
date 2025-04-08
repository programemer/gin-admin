package response

import "github.com/programemer/gin-admin/model/entity"

type SysAPIResponse struct {
	Api entity.SysApi `json:"api"`
}

type SysAPIListResponse struct {
	Apis []entity.SysApi `json:"apis"`
}

type SysSyncApis struct {
	NewApis    []entity.SysApi `json:"newApis"`
	DeleteApis []entity.SysApi `json:"deleteApis"`
}
