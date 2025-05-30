package dto

type MLBotInfo struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	UID      string `json:"uid"`
	Protocol string `json:"protocol"`
}

type MLBotUpdate struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	UID      string `json:"uid"`
	Protocol string `json:"protocol" validate:"required,oneof=http https"`
}

type MLBotTest struct {
	Message string `json:"message" validate:"required"`
}
