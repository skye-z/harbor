package service

import (
	"harbor/model"
	"harbor/util"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type UserService struct {
	UserModel model.UserModel
}

func NewUserService(engine *xorm.Engine) *UserService {
	us := new(UserService)
	us.UserModel = model.UserModel{
		DB: engine,
	}
	return us
}

type loginRequest struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type loginResponse struct {
	User   *model.User `json:"user"`
	Token  string      `json:"token"`
	Expire int64       `json:"expire"`
	Time   int64       `json:"time"`
}

// 登录
func (us UserService) Login(ctx *gin.Context) {
	var form loginRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		util.ReturnMessage(ctx, false, "传入数据无效")
		return
	}
	if len(form.Name) == 0 || len(form.Pass) == 0 {
		util.ReturnMessage(ctx, false, "用户名或密码不能为空")
		return
	}
	if len(form.Pass) != 32 {
		util.ReturnMessage(ctx, false, "禁止传输明文密码")
		return
	}
	user, token, exp, err := us.GetUserLoginInfo(form.Name, form.Pass)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取用户信息出错")
		return
	}
	if user == nil {
		util.ReturnMessage(ctx, false, "用户名或密码错误")
		return
	}
	ctx.JSON(200, loginResponse{User: user, Token: token, Expire: exp, Time: time.Now().Unix()})
}

func (us UserService) GetUserLoginInfo(name string, pass string) (*model.User, string, int64, error) {
	users, err := us.UserModel.GetUserList(name, pass, 0, 1)
	if err != nil {
		return nil, "", 0, err
	}
	if len(users) != 1 {
		return nil, "", 0, nil
	}
	user := users[0]
	token, exp, err := GenerateToken(&user)
	if err != nil {
		return nil, "", 0, err
	}

	return &user, token, exp, nil
}
