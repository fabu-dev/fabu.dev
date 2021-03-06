package filter

import (
	"fabu.dev/api/application/controller/response"
	"fabu.dev/api/application/model"
	"fabu.dev/api/application/service"
	"fabu.dev/api/pkg/api"
	"fabu.dev/api/pkg/api/code"
	"fabu.dev/api/pkg/api/global"
	"fabu.dev/api/pkg/api/request"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type App struct {
	BaseFilter
	service *service.App
}

func NewApp() *App {
	return &App{
		service: service.NewApp(),
	}
}

// 获取一个团队的APP列表
func (f *App) GetList(c *gin.Context) (*response.AppList, *api.Error) {
	params := &request.AppIndexParams{}
	if err := c.ShouldBind(params); err != nil {
		return nil, api.NewError(code.ErrorRequest, err.Error())
	}

	logrus.Info("getList params", params)
	result, err := f.service.GetListByTeamId(params)

	return result, err
}

func (f *App) GetSquareList(c *gin.Context) (*response.AppList, *api.Error) {
	return f.service.GetPublicApps()
}

// 上传文件
func (f *App) Upload(c *gin.Context) *api.Error {
	params := &request.UploadParams{}
	if err := c.ShouldBind(params); err != nil {
		return api.NewError(code.ErrorRequest, err.Error())
	}

	operator := f.GetOperator(c)

	err := f.service.Upload(params, operator)

	return err
}

func (f *App) UploadByAPI(c *gin.Context) *api.Error {
	params := &request.UploadByAPIParams{}
	if err := c.ShouldBind(params); err != nil {
		return api.NewError(code.ErrorRequest, err.Error())
	}

	err := f.service.UploadByAPI(params)

	return err
}

// 解析上传好的app文件
func (f *App) GetAppInfoByIdentifier(c *gin.Context) (*global.AppInfo, *api.Error) {
	params := &request.AppInfoParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		logrus.Error("add member err : ", err)
		return nil, api.NewError(code.ErrorRequest, err.Error())
	}

	appInfo, err := f.service.GetAppInfoByIdentifier(params.Identifier)

	return appInfo, err
}

// 保存上传文件信息
func (f *App) Save(c *gin.Context) (*global.AppInfo, *api.Error) {
	params := &request.SaveParams{}
	if err := c.ShouldBindJSON(params); err != nil {
		return nil, api.NewError(code.ErrorRequest, err.Error())
	}

	operator := f.GetOperator(c)

	appInfo, err := f.service.Save(params, operator)

	return appInfo, err
}

// 验证获取APP详情 根据short_url
func (f *App) ViewByShort(c *gin.Context) (*model.AppInfo, *api.Error) {
	params := &request.AppViewByShortParams{}

	if err := c.ShouldBind(params); err != nil {
		return nil, api.NewError(code.ErrorRequest, err.Error())
	}

	app, err := f.service.GetInfoByShortUrl(params.ShortUrl)

	return app, err
}

// 验证获取APP详情
func (f *App) View(c *gin.Context) (*model.AppInfo, *api.Error) {
	params := &request.AppViewParams{}

	if err := c.ShouldBindUri(params); err != nil {
		return nil, api.NewError(code.ErrorRequest, err.Error())
	}

	// 调用service对应的方法
	app, err := f.service.GetInfoById(params.Id)

	return app, err
}

func (f *App) Delete(c *gin.Context) *api.Error {
	params := &request.AppDeleteParams{}

	if err := c.ShouldBindJSON(params); err != nil {
		return api.NewError(code.ErrorRequest, err.Error())
	}

	operator := f.GetOperator(c)

	// 调用service对应的方法
	err := f.service.Delete(params, operator)

	return err
}
