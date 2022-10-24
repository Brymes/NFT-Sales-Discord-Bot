package models

import (
	"DIA-NFT-Sales-Bot/config"
	"database/sql"
	"gorm.io/gorm"
)

type Subscriptions struct {
	gorm.Model
	Command    string         `gorm:"column:command;not null"`
	Blockchain string         `gorm:"column:blockchain;not null"`
	ChannelID  sql.NullString `gorm:"column:channel_id"`
	Address    sql.NullString `gorm:"column:address"`
	Threshold  float64        `gorm:"column:threshold"`
	All        bool           `gorm:"column:all"`
	Active     bool           `gorm:"column:is_active"`
}

func (subscription Subscriptions) SaveSubscription() {
	result := config.DBClient.Model(&subscription).Create(&subscription)

	if result.Error != nil {
		err := "Error Saving Subscription: \n" + result.Error.Error()
		panic(err)
	}
}

func (subscription Subscriptions) LoadChannelSubscriptions() []Subscriptions {
	var subscriptions []Subscriptions

	result := config.DBClient.Model(&subscription).Where("is_active = ? AND channel_id = ?", true, subscription.ChannelID.String).Find(&subscriptions)
	if result.Error != nil {
		err := "Error Fetching Subscriptions by Channel: \n" + result.Error.Error()
		panic(err)
	}

	return subscriptions
}

func (subscription Subscriptions) DeactivateChannelSubscriptions() {
	result := config.DBClient.Model(&subscription).Where("is_active = ? AND channel_id = ?", true, subscription.ChannelID.String).Update("is_active", false)
	if result.Error != nil {
		err := "Error Deactivating Channel Subscriptions : \n" + result.Error.Error()
		panic(err)
	}
}

func (subscription Subscriptions) DeactivateAllSubscriptions() {
	result := config.DBClient.Model(&subscription).Where("is_active = ?", true).Update("is_active", false)
	if result.Error != nil {
		err := "Error Deactivating All Subscriptions : \n" + result.Error.Error()
		panic(err)
	}
}

func (subscription Subscriptions) DeactivateSubscriptions(idList []int) {
	result := config.DBClient.Model(&subscription).Where("id IN ?", idList).Updates(map[string]interface{}{"is_active": false})
	if result.Error != nil {
		err := "Error batch deactivating All Subscriptions : \n" + result.Error.Error()
		panic(err)
	}
}

func (subscription Subscriptions) UnsubscribeChannelSalesUpdates() {
	var result *gorm.DB

	if subscription.Address.Valid == false {
		result = config.DBClient.Model(&subscription).Where("command = ? AND channel_id = ? AND threshold = ?", subscription.Command, subscription.ChannelID.String, subscription.Threshold).Update("is_active", false)
	} else {
		result = config.DBClient.Model(&subscription).Where("command = ? AND channel_id = ? AND address = ?", subscription.Command, subscription.ChannelID.String, subscription.Address.String).Update("is_active", false)
	}
	if result.Error != nil {
		err := "Error Deactivating Channel Sales Subscriptions : \n" + result.Error.Error()
		panic(err)
	}
}

func (subscription Subscriptions) UnsubscribeSalesUpdates() {

	result := config.DBClient.Model(&subscription).Where("command = ?", subscription.Command).Update("is_active", false)
	if result.Error != nil {
		err := "Error Deactivating Channel Subscriptions : \n" + result.Error.Error()
		panic(err)
	}
}

func (subscription Subscriptions) LoadAllSubscriptions() []Subscriptions {
	var subscriptions []Subscriptions

	result := config.DBClient.Model(Subscriptions{}).Where("is_active = ?", true).Find(&subscriptions)
	if result.Error != nil {
		err := "Error Fetching Subscriptions by Channel: \n" + result.Error.Error()
		panic(err)
	}

	return subscriptions
}
