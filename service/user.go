/*
用户服务

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package service

import (
	"strconv"
	"time"

	"github.com/skye-z/harbor/model"
	"github.com/skye-z/harbor/util"

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
		util.ReturnMessage(ctx, false, DATA_ERROR)
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
	logger := &model.LogModel{
		DB: us.UserModel.DB,
	}
	logger.AddLog("platform", "password", user.Nickname+" 从 "+ctx.ClientIP())
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

func (us UserService) State(ctx *gin.Context) {
	util.ReturnMessage(ctx, util.GetBool("oauth2.enable"), "")
}

// 获取用户列表
func (us UserService) GetList(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}
	list, err := us.UserModel.GetUserList("", "", 1, 10000)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取用户列表失败")
		return
	}
	util.ReturnData(ctx, true, list)
}

// 获取用户详情
func (us UserService) GetInfo(ctx *gin.Context) {
	var uid int64
	var err error
	id := ctx.Query("uid")
	if id != "" {
		uid, err = strconv.ParseInt(id, 10, 64)
	}
	if uid == 0 || err != nil {
		obj, exist := ctx.Get("user")
		if exist {
			user := obj.(model.User)
			uid = user.Id
		} else {
			util.ReturnMessage(ctx, false, "未指定要查询的用户")
			return
		}
	}
	if uid == 0 {
		util.ReturnMessage(ctx, false, "未指定要查询的用户")
		return
	}
	user := &model.User{
		Id: uid,
	}
	err = us.UserModel.GetUser(user)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取用户信息失败")
	} else {
		util.ReturnData(ctx, true, user)
	}
}

type FormUser struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Name     string `json:"name"`
	Admin    bool   `json:"admin"`
	Pass     string `json:"pass"`
}

// 添加用户
func (us UserService) Add(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}

	var addObj FormUser
	if err := ctx.ShouldBindJSON(&addObj); err != nil {
		util.ReturnMessage(ctx, false, DATA_ERROR)
		return
	}
	var form model.User = model.User{
		Nickname: addObj.Nickname,
		Name:     addObj.Name,
		Admin:    addObj.Admin,
		Pass:     addObj.Pass,
	}
	if len(form.Name) == 0 {
		util.ReturnMessage(ctx, false, "用户名不能为空")
		return
	}
	user := model.User{
		Name: form.Name,
	}
	us.UserModel.GetUser(&user)
	if user.Id != 0 {
		util.ReturnMessage(ctx, false, "用户已存在")
		return
	}
	if len(form.Nickname) == 0 {
		form.Nickname = form.Name
	}
	if len(form.Pass) == 0 {
		util.ReturnMessage(ctx, false, "密码不能为空")
		return
	}
	state := us.UserModel.AddUser(&form)
	if state {
		util.ReturnMessage(ctx, true, "用户添加成功")
	} else {
		util.ReturnMessage(ctx, false, "用户添加失败")
	}
}

// 删除用户
func (us UserService) Del(ctx *gin.Context) {
	if !util.CheckAuth(ctx) {
		util.ReturnMessage(ctx, false, "权限不足")
		return
	}

	id := ctx.Query("uid")
	if id == "" {
		util.ReturnMessage(ctx, false, "请先指定用户")
		return
	}
	uid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.ReturnMessage(ctx, false, "请先指定用户")
		return
	}

	if uid == 1 {
		util.ReturnMessage(ctx, false, "禁止删除1号用户")
		return
	}

	state := us.UserModel.DelUser(&model.User{
		Id: uid,
	})
	if state {
		util.ReturnMessage(ctx, true, "用户删除成功")
	} else {
		util.ReturnMessage(ctx, false, "用户删除失败")
	}
}

// 编辑用户
func (us UserService) Edit(ctx *gin.Context) {
	var addObj FormUser
	if err := ctx.ShouldBindJSON(&addObj); err != nil {
		util.ReturnMessage(ctx, false, DATA_ERROR)
		return
	}

	obj, exist := ctx.Get("user")
	if !exist {
		util.ReturnMessage(ctx, false, "请先登陆")
		return
	}
	user := obj.(model.User)
	if addObj.Id != user.Id && !user.Admin {
		util.ReturnMessage(ctx, false, "权限不足")
	}

	var form model.User = model.User{
		Nickname: addObj.Nickname,
		Name:     addObj.Name,
		Admin:    addObj.Admin,
		Pass:     addObj.Pass,
	}

	if form.Name != "" {
		user = model.User{
			Name: form.Name,
		}
		us.UserModel.GetUser(&user)
		if user.Id != 0 && user.Id != addObj.Id {
			util.ReturnMessage(ctx, false, "用户名已存在")
			return
		}
		if len(form.Nickname) == 0 {
			form.Nickname = form.Name
		}
	}

	form.Id = addObj.Id
	state := us.UserModel.EditUser(&form)
	if state {
		util.ReturnMessage(ctx, true, "用户编辑成功")
	} else {
		util.ReturnMessage(ctx, false, "用户编辑失败")
	}
}

// OAuth2绑定
func (us UserService) Bind(ctx *gin.Context) {
	id := ctx.Query("id")
	name := ctx.Query("name")
	obj, exist := ctx.Get("user")
	if !exist {
		util.ReturnMessage(ctx, false, "请先登陆")
		return
	}
	user := obj.(model.User)

	var form model.User = model.User{
		Id:        user.Id,
		OAuthId:   id,
		OAuthName: name,
	}
	state := us.UserModel.EditUser(&form)
	if state {
		util.ReturnMessage(ctx, true, "用户绑定成功")
	} else {
		util.ReturnMessage(ctx, false, "用户绑定失败")
	}
}
