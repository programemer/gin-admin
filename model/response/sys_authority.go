package response

import "github.com/programemer/gin-admin/model/entity"

type SysAuthorityResponse struct {
	Authority entity.SysAuthority `json:"authority"`
}

type SysAuthorityCopyResponse struct {
	Authority      entity.SysAuthority `json:"authority"`
	OldAuthorityId uint                `json:"oldAuthorityId"` // 旧角色ID
}
