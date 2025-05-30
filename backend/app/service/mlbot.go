package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/1Panel-dev/1Panel/backend/app/dto"
	"github.com/1Panel-dev/1Panel/backend/app/model"
	"github.com/1Panel-dev/1Panel/backend/app/repo"
	"github.com/1Panel-dev/1Panel/backend/constant"
	"github.com/1Panel-dev/1Panel/backend/global"
)

type MLBotService struct{}

type IMLBotService interface {
	GetConfig() (model.MLBotConfig, error)
	UpdateConfig(req dto.MLBotUpdate) error
	SendNotification(message string) error
	SendCronjobNotification(cronjob model.Cronjob, status string, errorMsg string) error
}

type SendMessageRequest struct {
	UID     *string `json:"uid"`
	Message string  `json:"message"`
}

var mlBotRepo = repo.NewIMLBotRepo()

func NewIMLBotService() IMLBotService {
	return &MLBotService{}
}

func (m *MLBotService) GetConfig() (model.MLBotConfig, error) {
	// 确保表存在（延迟初始化）
	m.ensureTablesExist()
	
	config, err := mlBotRepo.Get()
	if err != nil {
		// 如果没有配置，创建并返回默认配置
		defaultConfig := model.MLBotConfig{
			Enabled:       false, // 默认关闭，用户需要手动启用
			Host:          "223.254.129.240",
			Port:          3000,
			UID:           "master",
			Protocol:      "http",
			MessagePrefix: "[1Panel]",
		}
		
		// 尝试保存默认配置到数据库
		if saveErr := mlBotRepo.Create(defaultConfig); saveErr != nil {
			global.LOG.Warnf("Failed to save default mlBot config: %v", saveErr)
		}
		
		return defaultConfig, nil
	}
	return config, nil
}

func (m *MLBotService) ensureTablesExist() {
	if global.DB == nil {
		return // 数据库还未初始化
	}
	
	// 自动创建 mlBot 配置表
	if err := global.DB.AutoMigrate(&model.MLBotConfig{}); err != nil {
		global.LOG.Errorf("Failed to create MLBotConfig table: %v", err)
	}
	
	// 自动为 cronjob 表添加通知字段（如果不存在的话）
	if err := global.DB.AutoMigrate(&model.Cronjob{}); err != nil {
		global.LOG.Errorf("Failed to update Cronjob table: %v", err)
	}
}

func (m *MLBotService) UpdateConfig(req dto.MLBotUpdate) error {
	// 确保表存在
	m.ensureTablesExist()
	
	existingConfig, err := mlBotRepo.Get()
	config := model.MLBotConfig{
		Enabled:       req.Enabled,
		Host:          req.Host,
		Port:          req.Port,
		UID:           req.UID,
		Protocol:      req.Protocol,
		MessagePrefix: req.MessagePrefix,
	}
	
	if err != nil {
		// 如果不存在，创建新的
		return mlBotRepo.Create(config)
	}
	// 更新现有配置
	config.ID = existingConfig.ID
	return mlBotRepo.Update(config)
}

func (m *MLBotService) SendNotification(message string) error {
	config, err := m.GetConfig()
	if err != nil {
		global.LOG.Errorf("Failed to get mlBot config: %v", err)
		return err
	}

	if !config.Enabled {
		global.LOG.Info("mlBot notification is disabled, skipping")
		return nil
	}

	// 添加消息前缀
	finalMessage := message
	if config.MessagePrefix != "" {
		finalMessage = config.MessagePrefix + " " + message
	}

	// 构建请求
	url := fmt.Sprintf("%s://%s:%d/api/message/send", config.Protocol, config.Host, config.Port)
	
	var uid *string
	if config.UID != "" {
		uid = &config.UID
	}
	
	reqBody := SendMessageRequest{
		UID:     uid,
		Message: finalMessage,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		global.LOG.Errorf("Failed to marshal request body: %v", err)
		return err
	}

	// 发送HTTP请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		global.LOG.Errorf("Failed to send notification to mlBot: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		global.LOG.Errorf("mlBot notification failed with status: %d", resp.StatusCode)
		return fmt.Errorf("notification failed with status: %d", resp.StatusCode)
	}

	global.LOG.Info("mlBot notification sent successfully")
	return nil
}

// 发送定时任务通知
func (m *MLBotService) SendCronjobNotification(cronjob model.Cronjob, status string, errorMsg string) error {
	config, err := m.GetConfig()
	if err != nil || !config.Enabled {
		return nil
	}

	var message string
	if status == constant.StatusSuccess {
		if !cronjob.NotifyOnSuccess {
			return nil
		}
		message = fmt.Sprintf("✅ 定时任务执行成功\n任务名称: %s\n任务类型: %s\n执行时间: %s", 
			cronjob.Name, cronjob.Type, time.Now().Format("2006-01-02 15:04:05"))
	} else if status == constant.StatusFailed {
		if !cronjob.NotifyOnFailure {
			return nil
		}
		message = fmt.Sprintf("❌ 定时任务执行失败\n任务名称: %s\n任务类型: %s\n错误信息: %s\n执行时间: %s", 
			cronjob.Name, cronjob.Type, errorMsg, time.Now().Format("2006-01-02 15:04:05"))
	}

	if message != "" {
		return m.SendNotification(message)
	}
	return nil
}
