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
		adminUser := CreateDefaultUser()
		userModel.AddUser(adminUser)
		util.Set("basic.install", "1")
		log.Println("[Data] init data")
	}
	log.Println("[Data] loading completed")
}

func CreateDefaultUser() *model.User {
	pwd := util.GenerateRandomString(12)
	log.Println("[Tips] create default admin, password is", pwd)
	return &model.User{
		Name:     "admin",
		Nickname: "管理员",
		Admin:    true,
		Pass:     util.CalculateMD5(pwd),
	}
}
