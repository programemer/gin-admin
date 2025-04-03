package system

import (
	"errors"
	"github.com/programemer/gin-admin/global"
	"github.com/programemer/gin-admin/model/entity"
	"gorm.io/gorm"
	"strings"
)

type ApiService struct {
}

var ApiServiceApp = new(ApiService)

func (apiService *ApiService) CreateApi(api entity.SysApi) (err error) {
	if !errors.Is(global.GVA_DB.Where("path = ? and  method = ?", api.Path, api.Method).First(&entity.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同api")
	}
	return global.GVA_DB.Create(&api).Error
}

func (apiService *ApiService) GetApiGroups(api entity.SysApi) (groups []string, groupApiMap map[string]string, err error) {
	var apis []entity.SysApi
	err = global.GVA_DB.Find(&apis).Error
	if err != nil {
		return
	}
	groupApiMap = make(map[string]string, 0)
	for i := range apis {
		pathArr := strings.Split(apis[i].Path, "/")
		newGroup := true
		for i2 := range groups {
			if groups[i2] == apis[i].ApiGroup {
				newGroup = false
			}
		}
		if newGroup {
			groups = append(groups, apis[i].ApiGroup)
		}
		groupApiMap[pathArr[1]] = apis[i].ApiGroup
	}
	return
}
