package data

import "time"

// 用户模型
type User struct {
	ID        int       `json:"id" xorm:"pk autoincr"`
	Username  string    `json:"username" xorm:"unique"`
	Password  string    `json:"password"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

// 系统日志模型
type SystemLog struct {
	ID        int64     `json:"id" xorm:"pk autoincr"`
	Type      string    `json:"type"`       // 日志类型: system/container/image/network/volume
	Level     string    `json:"level"`      // 日志级别: info/warning/error
	Action    string    `json:"action"`     // 操作类型
	Target    string    `json:"target"`     // 操作目标
	TargetID  string    `json:"target_id"`  // 目标ID
	Message   string    `json:"message"`    // 日志消息
	UserID    int       `json:"user_id"`    // 操作用户ID
	Username  string    `json:"username"`   // 操作用户名
	IPAddress string    `json:"ip_address"` // 客户端IP
	CreatedAt time.Time `json:"created_at"` // 创建时间
}
