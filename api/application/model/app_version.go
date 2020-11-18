package model

import (
	"errors"

	"fabu.dev/api/pkg/constant"

	"fabu.dev/api/pkg/api"
	"fabu.dev/api/pkg/utils"
	"github.com/jinzhu/gorm"
)

type AppVersion struct {
	DetailColumns []string
	BaseModel
}

type AppVersionInfo struct {
	Id          uint64         `json:"id" gorm:"primary_key"`
	AppId       uint64         `json:"app_id"`
	Tag         string         `json:"tag"`
	Code        string         `json:"code"`
	Description string         `json:"description"`
	Size        uint64         `json:"size"`
	Hash        string         `json:"hash"`
	Path        string         `json:"path"`
	IsPublish   uint64         `json:"is_publish"`
	ShortUrl    string         `json:"short_url"`
	QrCode      string         `json:"qr_code"`
	Status      uint8          `json:"status"`
	CreatedBy   string         `json:"created_by"`
	CreatedAt   utils.JSONTime `json:"created_at" gorm:"-"` // 插入时忽略该字段
	UpdatedBy   string         `json:"updated_by"`
	UpdatedAt   string         `json:"updated_at" gorm:"-"` // 插入时忽略该字段
}

func NewAppVersion() *AppVersion {
	AppVersion := &AppVersion{
		DetailColumns: []string{"id", "app_id", "tag", "code", "description", "size", "hash", "path", "is_publish", "short_url", "qr_code", "status", "updated_at", "updated_by", "created_at", "created_by"},
	}

	AppVersion.SetTableName("app_version")

	return AppVersion
}

// 添加app版本信息
func (m *AppVersion) Add(appVersion *AppVersionInfo) *api.Error {
	err := m.Db().Create(appVersion).Error

	return m.ProcessError(err)
}

// 编辑app信息
func (m *AppVersion) Edit(appVersion *AppVersionInfo) *api.Error {
	err := m.Db().Where("id = ?", appVersion.Id).Updates(appVersion).Error

	return m.ProcessError(err)
}

// 通过版本号获取版本信息
func (m *AppVersion) GetInfoByCode(AppId uint64, code string) (*AppVersionInfo, *api.Error) {
	appInfo := &AppVersionInfo{}
	err := m.Db().Select(m.DetailColumns).Where("app_id = ? and code = ?", AppId, code).First(appInfo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return appInfo, m.ProcessError(err)
}

// 获取一个app的版本列表
func (m *AppVersion) GetListByAppId(AppId uint64) ([]*AppVersionInfo, *api.Error) {
	appSlice := make([]*AppVersionInfo, 0, 32)
	err := m.Db().Select(m.DetailColumns).Where("app_id = ? and status = ?", AppId, constant.StatusEnable).Find(&appSlice).Order("code desc").Error

	return appSlice, m.ProcessError(err)
}

// 通过版本号获取版本信息
func (m *AppVersion) GetInfoByShortKey(key string) (*AppVersionInfo, *api.Error) {
	appInfo := &AppVersionInfo{}
	err := m.Db().Select(m.DetailColumns).Where("short_url = ?", key).First(appInfo).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return appInfo, m.ProcessError(err)
}
