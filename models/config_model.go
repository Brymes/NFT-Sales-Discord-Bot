package models

import "DIA-NFT-Sales-Bot/config"

type ConfigModel struct {
	TrackerCurrency string `gorm:"column:currency;not null"`
}

func (cm *ConfigModel) GetConfig() {
	result := config.DBClient.Model(&ConfigModel{}).Last(cm)

	if result.Error != nil {
		config.TrackerCurrency = "BrOkEn"
	} else {
		config.TrackerCurrency = cm.TrackerCurrency
	}
}

func (cm *ConfigModel) SaveConfig() {
	result := config.DBClient.Model(&ConfigModel{}).Create(cm)

	if result.Error != nil {
		err := "Error Saving Config Model: \n" + result.Error.Error()
		panic(err)
	}
}
