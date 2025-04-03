package system

import (
	"context"
	"github.com/pkg/errors"
	"github.com/programemer/gin-admin/global"
	common "github.com/programemer/gin-admin/model/common/request"
	"github.com/programemer/gin-admin/model/entity"
	"github.com/programemer/gin-admin/model/request"
)

type autoCodeHistory struct{}

// Create 创建代码生成器历史记录
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [songzhibin97](https://github.com/songzhibin97)
func (s *autoCodeHistory) Create(ctx context.Context, info request.SysAutoHistoryCreate) error {
	create := info.Create()
	err := global.GVA_DB.WithContext(ctx).Create(&create).Error
	if err != nil {
		return errors.Wrap(err, "创建失败")
	}
	return nil
}

// First 根据id获取代码生成器历史的数据
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [songzhibin97](https://github.com/songzhibin97)
func (s *autoCodeHistory) First(ctx context.Context, info common.GetById) (string, error) {
	var meta string
	err := global.GVA_DB.WithContext(ctx).Model(entity.SysAutoCodeHistory{}).Where("id = ?", info.ID).Pluck("request", &meta).Error
	if err != nil {
		return "", errors.Wrap(err, "获取失败!")
	}
	return meta, nil
}

// Repeat 检测重复
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [songzhibin97](https://github.com/songzhibin97)
func (s *autoCodeHistory) Repeat(businessDB, structName, abbreviation, Package string) bool {
	var count int64
	global.GVA_DB.Model(&entity.SysAutoCodeHistory{}).Where("business_db = ? and (struct_name = ? OR abbreviation = ?) and package = ? and flag = ?", businessDB, structName, abbreviation, Package, 0).Count(&count).Debug()
	return count > 0
}

// RollBack 回滚
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [songzhibin97](https://github.com/songzhibin97)
func (s *autoCodeHistory) RollBack(ctx context.Context, info request.SysAutoHistoryRollBack) error {
	var history entity.SysAutoCodeHistory
	err := global.GVA_DB.Where("id = ?", info.ID).First(&history).Error
	if err != nil {
		return err
	}
	if history.ExportTemplateID != 0 {
		err = global.GVA_DB.Delete(&entity.SysExportTemplate{}, "id = ?", history.ExportTemplateID).Error
		if err != nil {
			return err
		}
	}
	if info.DeleteApi {
		ids := info.ApiIds(history)
		err = ApiSer
	}
}
