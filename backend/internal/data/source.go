package data

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

func InitDB() *xorm.Engine {
	engine, err := xorm.NewEngine("sqlite", "./harbor.db")
	if err != nil {
		panic(err)
	}
	return engine
}

func InitDBTable(engine *xorm.Engine) {
	err := engine.Sync2(new(User))
	if err != nil {
		panic(err)
	}

	count, err := engine.Count(&User{})
	if err != nil {
		panic(err)
	}
	if count == 0 {
		// 默认密码 "harbor-skye"
		defaultPasswordMD5 := "58c823f11203f20fda6b4deb81d30b3b"

		if len(defaultPasswordMD5) != 32 {
			panic("Default password must be a 32-character MD5 hash")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPasswordMD5), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		admin := User{
			Username:  "admin",
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
		}
		_, err = engine.Insert(&admin)
		if err != nil {
			panic(err)
		}
	}
}

var Engine *xorm.Engine

// 获取数据库引擎
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
