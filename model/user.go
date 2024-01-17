package model

import (
	"fmt"

	"xorm.io/xorm"
)

type User struct {
	Id        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	OAuthId   string `json:"-"`
	OAuthName string `json:"oauthName"`
	Name      string `json:"name"`
	Admin     bool   `json:"admin"`
	Pass      string `json:"-"`
}

type UserModel struct {
	DB *xorm.Engine
}

// 获取用户
func (model UserModel) GetUser(user *User) error {
	has, err := model.DB.Get(user)
	if !has {
		return err
	}
	return nil
}

// 获取授权用户
func (model UserModel) GetOAuthUser(oauthId string) (*User, error) {
	user := &User{
		OAuthId: oauthId,
	}
	has, err := model.DB.Get(user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return user, nil
}

// 获取用户列表
func (model UserModel) GetUserList(name string, pass string, page int, num int) ([]User, error) {
	var users []User
	var err error
	if len(name) == 0 {
		err = model.DB.Limit(page*num, (page-1)*num).Find(&users)
	} else if len(pass) == 32 {
		err = model.DB.Where("name = ? and pass = ?", name, pass).Limit(1, 0).Find(&users)
	} else {
		err = model.DB.Where("name like ?", fmt.Sprintf("%%%s%%", name)).Limit(page*num, (page-1)*num).Find(&users)
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 添加用户
func (model UserModel) AddUser(user *User) bool {
	_, err := model.DB.Insert(user)
	return err == nil
}

// 编辑用户
func (model UserModel) EditUser(user *User) bool {
	if user.Id == 0 {
		return false
	}
	_, err := model.DB.ID(user.Id).Update(user)
	return err == nil
}

// 删除用户
func (model UserModel) DelUser(user *User) bool {
	if user.Id == 0 {
		return false
	}
	_, err := model.DB.Delete(user)
	return err == nil
}
