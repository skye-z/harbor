package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

// 全局数据库引擎
var Engine *xorm.Engine

// 初始化数据库连接
func InitDB() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite", "./harbor.db")
	if err != nil {
		return nil, errors.New("创建数据库引擎失败: " + err.Error())
	}
	Engine = engine
	return engine, nil
}

// 初始化数据库表
func InitDBTable(engine *xorm.Engine) error {
	err := engine.Sync2(new(User))
	if err != nil {
		return errors.New("同步数据库表失败: " + err.Error())
	}

	err = engine.Sync2(new(SystemLog))
	if err != nil {
		return errors.New("同步日志表失败: " + err.Error())
	}

	count, err := engine.Count(&User{})
	if err != nil {
		return errors.New("统计用户数量失败: " + err.Error())
	}
	if count == 0 {
		defaultPassword := "HarborAdmin2026!"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("加密默认密码失败: " + err.Error())
		}
		admin := User{
			Username:  "admin",
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
		}
		_, err = engine.Insert(&admin)
		if err != nil {
			return errors.New("创建默认管理员用户失败: " + err.Error())
		}
	}
	return nil
}

// 获取数据库引擎实例
func GetDB() *xorm.Engine {
	return Engine
}

// 关闭数据库连接
func Close() error {
	if Engine != nil {
		return Engine.Close()
	}
	return nil
}
