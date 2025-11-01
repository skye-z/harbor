package data

import "time"

// 用户模型
type User struct {
	ID        int       `json:"id" xorm:"pk autoincr"`  // 用户 ID，主键，自增
	Username  string    `json:"username" xorm:"unique"` // 用户名，唯一
	Password  string    `json:"password"`               // 密码
	IsAdmin   bool      `json:"is_admin"`               // 是否管理员
	CreatedAt time.Time `json:"created_at"`             // 创建时间
}
