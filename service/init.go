package service

import (
	"harbor/model"
	"harbor/util"
	"log"

	"xorm.io/xorm"
)

func InitDatabase(engine *xorm.Engine) {
	log.Println("[Data] load data")
	err := engine.Sync2(new(model.User))
	if err != nil {
		panic(err)
	}
	if !util.GetBool("basic.install") {
		userModel := model.UserModel{DB: engine}
		adminUser := model.User{
			Name:     "admin",
			Nickname: "管理员",
			Admin:    true,
			Pass:     "21232f297a57a5a743894a0e4a801fc3",
		}
		userModel.AddUser(&adminUser)
		util.Set("basic.install", "1")
		log.Println("[Data] init data")
	}
	log.Println("[Data] loading completed")
}
