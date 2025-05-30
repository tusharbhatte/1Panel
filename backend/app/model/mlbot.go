package model

type MLBotConfig struct {
	BaseModel

	Enabled       bool   `gorm:"type:tinyint(1);default:1" json:"enabled"`
	Host          string `gorm:"type:varchar(255);default:'223.254.129.240'" json:"host"`
	Port          int    `gorm:"type:int;default:3000" json:"port"`
	UID           string `gorm:"type:varchar(64)" json:"uid"`
	Protocol      string `gorm:"type:varchar(10);default:'http'" json:"protocol"`
	MessagePrefix string `gorm:"type:varchar(128);default:'[1Panel]'" json:"messagePrefix"`
}
