package v1

import (
	"github.com/1Panel-dev/1Panel/backend/app/api/v1/helper"
	"github.com/1Panel-dev/1Panel/backend/app/dto"
	"github.com/1Panel-dev/1Panel/backend/app/service"
	"github.com/1Panel-dev/1Panel/backend/constant"
	"github.com/gin-gonic/gin"
)

var mlBotService = service.NewIMLBotService()

// @Tags MLBot
// @Summary Get mlBot configuration
// @Accept json
// @Success 200 {object} dto.MLBotInfo
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /mlbot/config [get]
func (b *BaseApi) GetMLBotConfig(c *gin.Context) {
	config, err := mlBotService.GetConfig()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	
	info := dto.MLBotInfo{
		Enabled:  config.Enabled,
		Host:     config.Host,
		Port:     config.Port,
		UID:      config.UID,
		Protocol: config.Protocol,
	}
	helper.SuccessWithData(c, info)
}

// @Tags MLBot
// @Summary Update mlBot configuration
// @Accept json
// @Param request body dto.MLBotUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /mlbot/config [put]
// @x-panel-log {"bodyKeys":["enabled"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"更新 mlBot 配置","formatEN":"update mlBot config"}
func (b *BaseApi) UpdateMLBotConfig(c *gin.Context) {
	var req dto.MLBotUpdate
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := mlBotService.UpdateConfig(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags MLBot
// @Summary Test mlBot notification
// @Accept json
// @Param request body dto.MLBotTest true "request"
// @Success 200
// @Security ApiKeyAuth
// @Security Timestamp
// @Router /mlbot/test [post]
func (b *BaseApi) TestMLBotNotification(c *gin.Context) {
	var req dto.MLBotTest
	if err := helper.CheckBindAndValidate(&req, c); err != nil {
		return
	}

	if err := mlBotService.SendNotification(req.Message); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
