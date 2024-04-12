/*
日志数据模型

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package model

import (
	"encoding/json"
	"log"
	"time"

	"xorm.io/xorm"
)

type Log struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	Timestamp int64  `json:"timestamp"`
}

type LogModel struct {
	DB *xorm.Engine
}

// 添加日志
func (model LogModel) AddLog(types string, action string, details interface{}) bool {
	detailsMsg := ""
	switch v := details.(type) {
	case string:
		detailsMsg = v
	default:
		data, err := json.Marshal(details)
		if err == nil {
			detailsMsg = string(data)
		}
	}
	logs := &Log{
		Type:    types,
		Action:  action,
		Details: detailsMsg,
	}
	currentTime := time.Now()
	logs.Timestamp = currentTime.UnixMilli()
	_, err := model.DB.Insert(logs)
	return err == nil
}

// 获取日志列表
func (model LogModel) GetLogs(page int, num int) ([]Log, error) {
	var logs []Log
	err := model.DB.Desc("id").Limit(page*num, (page-1)*num).Find(&logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// 获取日志数量
func (model LogModel) GetLogNumber() (int64, error) {
	var logs Log
	num, err := model.DB.Count(logs)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// 清空日志
func (model LogModel) Clean() bool {
	_, err := model.DB.Exec("DELETE FROM log")
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
