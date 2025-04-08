package system

import (
	"errors"
	"github.com/programemer/gin-admin/global"
	"github.com/programemer/gin-admin/model/common/request"
	"github.com/programemer/gin-admin/model/entity"
	systemReq "github.com/programemer/gin-admin/model/request"
	"github.com/programemer/gin-admin/model/response"
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

// @function: CopyAuthority
// @description: 复制一个角色
// @param: copyInfo response.SysAuthorityCopyResponse
// @return: authority system.SysAuthority, err error
func (authorityService *AuthorityService) CopyAuthority(adminAuthorityID uint, copyInfo response.SysAuthorityCopyResponse) (authority entity.SysAuthority, err error) {
	var authorityBox entity.SysAuthority
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", copyInfo.Authority.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return authority, ErrRoleExistence
	}
	copyInfo.Authority.Children = []entity.SysAuthority{}
	menus, err := MenuServiceApp.GetMenuAuthority(&request.GetAuthorityId{AuthorityId: copyInfo.OldAuthorityId})
	if err != nil {
		return
	}
	var baseMenu []entity.SysBaseMenu
	for _, v := range menus {
		intNum := v.MenuId
		v.SysBaseMenu.ID = uint(intNum)
		baseMenu = append(baseMenu, v.SysBaseMenu)
	}
	copyInfo.Authority.SysBaseMenus = baseMenu
	err = global.GVA_DB.Create(&copyInfo.Authority).Error
	if err != nil {
		return
	}

	var btns []entity.SysAuthorityBtn

	err = global.GVA_DB.Find(&btns, "authority_id = ?", copyInfo.OldAuthorityId).Error
	if err != nil {
		return
	}
	if len(btns) > 0 {
		for i := range btns {
			btns[i].AuthorityId = copyInfo.Authority.AuthorityId
		}
		err = global.GVA_DB.Create(&btns).Error

		if err != nil {
			return
		}
	}
	paths := CasbinServiceApp.GetPolicyPathByAuthorityId(copyInfo.OldAuthorityId)
	err = CasbinServiceApp.UpdateCasbin(adminAuthorityID, copyInfo.Authority.AuthorityId, paths)
	if err != nil {
		_ = authorityService.DeleteAuthority(&copyInfo.Authority)
	}
	return copyInfo.Authority, err
}

//@function: UpdateAuthority
//@description: 更改一个角色
//@param: auth model.SysAuthority
//@return: authority system.SysAuthority, err error

func (authorityService *AuthorityService) UpdateAuthority(auth entity.SysAuthority) (authority entity.SysAuthority, err error) {
	var oldAuthority entity.SysAuthority
	err = global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&oldAuthority).Error
	if err != nil {
		global.GVA_LOG.Debug(err.Error())
		return entity.SysAuthority{}, errors.New("查询角色数据失败")
	}
	err = global.GVA_DB.Model(&oldAuthority).Updates(&auth).Error
	return auth, err
}

//@function: DeleteAuthority
//@description: 删除角色
//@param: auth *model.SysAuthority
//@return: err error

func (authorityService *AuthorityService) DeleteAuthority(auth *entity.SysAuthority) error {
	if errors.Is(global.GVA_DB.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色不存在")
	}
	if len(auth.Users) != 0 {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&entity.SysUser{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色有用户正在使用禁止删除")
	}
	if !errors.Is(global.GVA_DB.Where("parent_id = ?", auth.AuthorityId).First(&entity.SysAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var err error
		if err = tx.Preload("SysBaseMenus").Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(auth).Unscoped().Delete(auth).Error; err != nil {
			return err
		}

		if len(auth.SysBaseMenus) > 0 {
			if err = tx.Model(auth).Association("SysBaseMenus").Delete(auth.SysBaseMenus); err != nil {
				return err
			}
			// err = db.Association("SysBaseMenus").Delete(&auth)
		}
		if len(auth.DataAuthorityId) > 0 {
			if err = tx.Model(auth).Association("DataAuthorityId").Delete(auth.DataAuthorityId); err != nil {
				return err
			}
		}

		if err = tx.Delete(&entity.SysUserAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error; err != nil {
			return err
		}
		if err = tx.Where("authority_id = ?", auth.AuthorityId).Delete(&[]entity.SysAuthorityBtn{}).Error; err != nil {
			return err
		}

		authorityId := strconv.Itoa(int(auth.AuthorityId))

		if err = CasbinServiceApp.RemoveFilteredPolicy(tx, authorityId); err != nil {
			return err
		}

		return nil
	})
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

//@function: SetMenuAuthority
//@description: 菜单与角色绑定
//@param: auth *model.SysAuthority
//@return: error

func (authorityService *AuthorityService) SetMenuAuthority(auth *entity.SysAuthority) error {
	var s entity.SysAuthority
	global.GVA_DB.Preload("SysBaseMenus").First(&s, "authority_id = ?", auth.AuthorityId)
	err := global.GVA_DB.Model(&s).Association("SysBaseMenus").Replace(&auth.SysBaseMenus)
	return err
}

//@function: findChildrenAuthority
//@description: 查询子角色
//@param: authority *model.SysAuthority
//@return: err error

func (authorityService *AuthorityService) findChildrenAuthority(authority *entity.SysAuthority) (err error) {
	err = global.GVA_DB.Preload("DataAuthorityId").Where("parent_id = ?", authority.AuthorityId).Find(&authority.Children).Error
	if len(authority.Children) > 0 {
		for k := range authority.Children {
			err = authorityService.findChildrenAuthority(&authority.Children[k])
		}
	}
	return err
}
