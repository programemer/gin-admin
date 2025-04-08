package system

import (
	"errors"
	"github.com/programemer/gin-admin/global"
	"github.com/programemer/gin-admin/model/entity"
	systemReq "github.com/programemer/gin-admin/model/request"
	"gorm.io/gorm"
	"strconv"
)

var ErrRoleExistence = errors.New("存在相同角色id")

// @author: [piexlmax](https://github.com/piexlmax)
// @function: CreateAuthority
// @description: 创建一个角色
// @param: auth model.SysAuthority
// @return: authority system.SysAuthority, err error
type AuthorityService struct{}

var AuthorityServiceApp = new(AuthorityService)

func (authorityService *AuthorityService) CreateAuthority(auth entity.SysAuthority) (authority entity.SysAuthority, err error) {
	if err = global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&auth).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		return auth, ErrRoleExistence
	}

	e := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err = tx.Create(&auth).Error; err != nil {
			return err
		}
		auth.SysBaseMenus = systemReq.DefaultMenu()
		if err = tx.Model(&auth).Association("SysBaseMenus").Replace(&auth.SysBaseMenus); err != nil {
			return err
		}

		casbinInfos := systemReq.DefaultCasbin()
		authorityId := strconv.Itoa(int(auth.AuthorityId))
		rules := [][]string{}
		for _, v := range casbinInfos {
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
		return CasbinServiceApp.AddPolicies(tx, rules)
	})
	return auth, e
}

// @author: liyu
// @function: GetAuthorityInfoList
// @description: 分页获取数据
// @param: info request.PageInfo
// @return: list interface{}, total int64, err error
func (authorityService *AuthorityService) GetStructAuthorityList(authorityID uint) (list []uint, err error) {
	var auth entity.SysAuthority
	_ = global.GVA_DB.First(&auth, "authority_id = ?", authorityID).Error
	var authorities []entity.SysAuthority
	err = global.GVA_DB.Preload("DataAuthorityId").Where("parent_id = ?", authorityID).Find(&authorities).Error
	if len(authorities) > 0 {
		for k := range authorities {
			list = append(list, authorities[k].AuthorityId)
			childrenList, err := authorityService.GetStructAuthorityList(authorities[k].AuthorityId)
			if err == nil {
				list = append(list, childrenList...)
			}
		}
	}
	if *auth.ParentId == 0 {
		list = append(list, authorityID)
	}
	return list, err
}

func (authorityService *AuthorityService) CheckAuthorityIDAuth(authorityID, targetID uint) (err error) {
	if !global.GVA_CONFIG.System.UseStrictAuth {
		return nil
	}
	authIds, err := authorityService.GetStructAuthorityList(authorityID)
	if err != nil {
		return err
	}
	hasAuth := false
	for _, v := range authIds {
		if v == targetID {
			hasAuth = true
			break
		}
	}

	if !hasAuth {
		return errors.New("您提交的角色ID不合法")
	}
	return nil
}

func (authorityService *AuthorityService) GetParentAuthorityID(authorityID uint) (parentID uint, err error) {
	var authority entity.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", authorityID).First(&authority).Error
	return *authority.ParentId, err
}
