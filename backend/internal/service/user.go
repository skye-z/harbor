package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/data"
	"github.com/skye-z/harbor/internal/util/response"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

// 用户服务
type UserService struct {
	engine *xorm.Engine
}

// 创建用户服务实例
func NewUserService(engine *xorm.Engine) *UserService {
	return &UserService{engine: engine}
}

// 获取用户列表
func (s *UserService) GetList(c *gin.Context) {
	var users []data.User
	err := s.engine.Find(&users)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, users)
}

// 创建用户请求结构
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin"`
}

// 创建用户
func (s *UserService) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser data.User
	has, err := s.engine.Where("username = ?", req.Username).Get(&existingUser)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	if has {
		response.BadRequest(c, "用户名已存在")
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, "密码加密失败")
		return
	}

	// 创建用户
	user := &data.User{
		Username: req.Username,
		Password: string(hashedPassword),
		IsAdmin:  req.IsAdmin,
	}
	_, err = s.engine.Insert(user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "用户创建成功", user)
}

// 更新用户请求结构
type UpdateUserRequest struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  *bool  `json:"is_admin"`
}

// 更新用户
func (s *UserService) Update(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 获取原用户信息
	var user data.User
	has, err := s.engine.ID(req.ID).Get(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	if !has {
		response.BadRequest(c, "用户不存在")
		return
	}

	// 检查是否是唯一管理员降级为普通用户
	if req.IsAdmin != nil && !*req.IsAdmin && user.IsAdmin {
		// 查询管理员数量
		adminCount, err := s.engine.Where("is_admin = ?", true).Count(&data.User{})
		if err != nil {
			response.Error(c, err.Error())
			return
		}
		if adminCount <= 1 {
			response.BadRequest(c, "系统必须保留至少一个管理员")
			return
		}
	}

	// 更新字段
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(c, "密码加密失败")
			return
		}
		user.Password = string(hashedPassword)
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	_, err = s.engine.ID(req.ID).Update(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "用户更新成功", user)
}

// 删除用户
func (s *UserService) Delete(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		response.BadRequest(c, "缺少用户ID")
		return
	}

	// 获取用户信息
	var user data.User
	has, err := s.engine.ID(id).Get(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	if !has {
		response.BadRequest(c, "用户不存在")
		return
	}

	// 检查是否是管理员
	if user.IsAdmin {
		// 查询管理员数量
		adminCount, err := s.engine.Where("is_admin = ?", true).Count(&data.User{})
		if err != nil {
			response.Error(c, err.Error())
			return
		}
		if adminCount <= 1 {
			response.BadRequest(c, "系统必须保留至少一个管理员")
			return
		}
	}

	_, err = s.engine.ID(id).Delete(&data.User{})
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "用户已删除", nil)
}

// 获取当前用户信息
func (s *UserService) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var user data.User
	has, err := s.engine.ID(userID).Get(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	if !has {
		response.BadRequest(c, "用户不存在")
		return
	}

	// 不返回密码
	user.Password = ""
	response.Success(c, user)
}

// 修改密码请求结构
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// 修改密码
func (s *UserService) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 验证新密码长度
	if len(req.NewPassword) < 6 {
		response.BadRequest(c, "新密码长度不能少于6位")
		return
	}

	// 获取用户信息
	var user data.User
	has, err := s.engine.ID(userID).Get(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	if !has {
		response.BadRequest(c, "用户不存在")
		return
	}

	// 验证旧密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		response.BadRequest(c, "旧密码错误")
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, "密码加密失败")
		return
	}

	// 更新密码
	user.Password = string(hashedPassword)
	_, err = s.engine.ID(userID).Update(&user)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "密码修改成功", nil)
}

// 错误定义
var (
	ErrUserNotFound     = errors.New("用户不存在")
	ErrUsernameExists   = errors.New("用户名已存在")
	ErrLastAdmin        = errors.New("系统必须保留至少一个管理员")
	ErrPasswordTooShort = errors.New("密码长度不能少于6位")
)
