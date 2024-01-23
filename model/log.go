package model

import (
	"encoding/json"
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
func (model LogModel) AddLog(_type string, action string, details interface{}) bool {
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
		Type:    _type,
		Action:  action,
		Details: detailsMsg,
	}
	currentTime := time.Now()
	logs.Timestamp = currentTime.UnixNano()
	_, err := model.DB.Insert(logs)
	return err == nil
}

// 获取日志列表
func (model LogModel) GetLogs(page int, num int) ([]Log, error) {
	var logs []Log
	err := model.DB.Limit(page*num, (page-1)*num).Find(&logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
