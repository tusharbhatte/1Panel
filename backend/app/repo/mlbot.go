package repo

import (
	"github.com/1Panel-dev/1Panel/backend/app/model"
	"github.com/1Panel-dev/1Panel/backend/global"
)

type MLBotRepo struct{}

type IMLBotRepo interface {
	Get() (model.MLBotConfig, error)
	Create(config model.MLBotConfig) error
	Update(config model.MLBotConfig) error
}

func NewIMLBotRepo() IMLBotRepo {
	return &MLBotRepo{}
}

func (m *MLBotRepo) Get() (model.MLBotConfig, error) {
	var config model.MLBotConfig
	err := global.DB.First(&config).Error
	return config, err
}

func (m *MLBotRepo) Create(config model.MLBotConfig) error {
	return global.DB.Create(&config).Error
}

func (m *MLBotRepo) Update(config model.MLBotConfig) error {
	return global.DB.Save(&config).Error
}
