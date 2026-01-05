package service

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/skye-z/harbor/internal/data"
	"github.com/skye-z/harbor/internal/util/config"
	"github.com/skye-z/harbor/internal/util/response"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

// 认证服务
type AuthService struct {
	engine *xorm.Engine
}

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Passkey  string `json:"passkey"`
}

// 登录响应结构
type LoginResponse struct {
	Token     string `json:"token"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	ExpiresAt int64  `json:"expires_at"`
}

// 通行密钥凭证
type PasskeyCredential struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	PublicKey string `json:"public_key"`
	Counter   int    `json:"counter"`
}

// 创建认证服务实例
func NewAuthService(engine *xorm.Engine) *AuthService {
	return &AuthService{engine: engine}
}

// 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user data.User
	has, err := s.engine.Where("username = ?", req.Username).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("用户不存在")
	}

	if req.Passkey != "" {
		return s.loginWithPasskey(req, &user)
	}

	return s.loginWithPassword(req, &user)
}

// 密码登录
func (s *AuthService) loginWithPassword(req *LoginRequest, user *data.User) (*LoginResponse, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return s.generateToken(user)
}

// 通行密钥登录
func (s *AuthService) loginWithPasskey(req *LoginRequest, user *data.User) (*LoginResponse, error) {
	var credentials []PasskeyCredential
	if err := s.engine.Where("user_id = ?", user.ID).Find(&credentials); err != nil {
		return nil, err
	}
	if len(credentials) == 0 {
		return nil, errors.New("未找到通行密钥")
	}

	signature := req.Passkey
	var matched *PasskeyCredential
	for i := range credentials {
		if verifyPasskeySignature(credentials[i].PublicKey, signature, req.Username) {
			matched = &credentials[i]
			break
		}
	}
	if matched == nil {
		return nil, errors.New("通行密钥验证失败")
	}

	return s.generateToken(user)
}

const DefaultTokenExpiration = 24 * time.Hour

// 生成认证令牌
func (s *AuthService) generateToken(user *data.User) (*LoginResponse, error) {
	secretKey := config.GetString("jwt.secret")
	if secretKey == "" {
		return nil, errors.New("JWT secret not configured")
	}

	expiration := DefaultTokenExpiration
	if expStr := config.GetString("jwt.expiration"); expStr != "" {
		if parsed, err := time.ParseDuration(expStr); err == nil {
			expiration = parsed
		}
	}

	expiresAt := time.Now().Add(expiration)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &LoginResponse{
		Token:     tokenString,
		UserID:    user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// 生成通行密钥对
func GeneratePasskeyPair() (publicKey, privateKey string, err error) {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	pubBytes, err := json.Marshal(private.PublicKey)
	if err != nil {
		return "", "", err
	}

	privBytes, err := json.Marshal(private)
	if err != nil {
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(pubBytes), base64.StdEncoding.EncodeToString(privBytes), nil
}

// 生成签名挑战（包含随机 nonce）
func GenerateChallenge() (string, string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", "", err
	}
	nonceStr := base64.StdEncoding.EncodeToString(nonce)
	challenge := base64.StdEncoding.EncodeToString([]byte(nonceStr + ":" + strconv.FormatInt(time.Now().Unix(), 10)))
	return challenge, nonceStr, nil
}

// 签名挑战
func SignChallenge(privateKey, challenge string) (string, error) {
	privBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", err
	}

	var private ecdsa.PrivateKey
	if err := json.Unmarshal(privBytes, &private); err != nil {
		return "", err
	}

	chalBytes := sha256.Sum256([]byte(challenge))
	r, s, err := ecdsa.Sign(rand.Reader, &private, chalBytes[:])
	if err != nil {
		return "", err
	}

	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sigBytes := append(rBytes, sBytes...)
	return base64.StdEncoding.EncodeToString(sigBytes), nil
}

// 验证签名挑战（包含 nonce 校验和过期时间）
func VerifyChallenge(nonce, challenge, signature, publicKey string) bool {
	pubBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false
	}

	var pub ecdsa.PublicKey
	if err := json.Unmarshal(pubBytes, &pub); err != nil {
		return false
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(sigBytes[:len(sigBytes)/2])
	s.SetBytes(sigBytes[len(sigBytes)/2:])

	chalBytes := sha256.Sum256([]byte(challenge))
	return ecdsa.Verify(&pub, chalBytes[:], r, s)
}

// 验证通行密钥签名
func verifyPasskeySignature(publicKey, signature, data string) bool {
	pubBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false
	}

	var pub ecdsa.PublicKey
	if err := json.Unmarshal(pubBytes, &pub); err != nil {
		return false
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	r := new(big.Int)
	s := new(big.Int)
	r.SetBytes(sigBytes[:len(sigBytes)/2])
	s.SetBytes(sigBytes[len(sigBytes)/2:])

	return ecdsa.Verify(&pub, []byte(data), r, s)
}

// 添加通行密钥
func (s *AuthService) AddPasskey(userID int, publicKey string) error {
	credential := &PasskeyCredential{
		UserID:    userID,
		PublicKey: publicKey,
		Counter:   0,
	}
	_, err := s.engine.Insert(credential)
	return err
}

// 删除通行密钥
func (s *AuthService) RemovePasskey(userID, credentialID string) error {
	id, _ := strconv.Atoi(credentialID)
	_, err := s.engine.Where("user_id = ? AND id = ?", userID, id).Delete(&PasskeyCredential{})
	return err
}

// 获取用户通行密钥列表
func (s *AuthService) GetPasskeys(userID int) ([]PasskeyCredential, error) {
	var credentials []PasskeyCredential
	err := s.engine.Where("user_id = ?", userID).Find(&credentials)
	return credentials, err
}

// 认证处理器
type AuthHandler struct {
}

// 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// 认证中间件
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			response.Unauthorized(c, "未授权访问")
			c.Abort()
			return
		}

		c.Next()
	}
}

// 令牌载荷
type TokenPayload struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

// 解析认证令牌
func ParseToken(tokenString string) (*TokenPayload, error) {
	secretKey := config.GetString("jwt.secret")
	if secretKey == "" {
		return nil, errors.New("JWT secret not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		payload := &TokenPayload{
			UserID:   int(claims["user_id"].(float64)),
			Username: claims["username"].(string),
			IsAdmin:  claims["is_admin"].(bool),
		}
		return payload, nil
	}

	return nil, errors.New("invalid token")
}

// 生成认证令牌
func GenerateToken(payload *TokenPayload, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  payload.UserID,
		"username": payload.Username,
		"is_admin": payload.IsAdmin,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
