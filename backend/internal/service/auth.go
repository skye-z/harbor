package service

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/data"
	"github.com/skye-z/harbor/internal/util/config"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

type AuthService struct {
	engine *xorm.Engine
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名，必填
	Password string `json:"password"`                    // 密码
	Passkey  string `json:"passkey"`                     // 通行密钥
}

type LoginResponse struct {
	Token     string `json:"token"`      // 认证令牌
	UserID    int    `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	IsAdmin   bool   `json:"is_admin"`   // 是否管理员
	ExpiresAt int64  `json:"expires_at"` // 过期时间戳
}

type PasskeyCredential struct {
	ID        string `json:"id"`         // 凭证ID
	PublicKey string `json:"public_key"` // 公钥
	UserID    int    `json:"user_id"`    // 关联用户ID
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

	// 根据请求参数选择登录方式
	if req.Passkey != "" {
		return s.loginWithPasskey(req, &user)
	}

	return s.loginWithPassword(req, &user)
}

// 使用密码登录
func (s *AuthService) loginWithPassword(req *LoginRequest, user *data.User) (*LoginResponse, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return s.generateToken(user)
}

// 使用通行密钥登录
func (s *AuthService) loginWithPasskey(req *LoginRequest, user *data.User) (*LoginResponse, error) {
	var credential PasskeyCredential
	has, err := s.engine.Where("user_id = ?", user.ID).Get(&credential)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("未找到通行密钥")
	}

	if credential.ID != req.Passkey {
		return nil, errors.New("通行密钥无效")
	}

	return s.generateToken(user)
}

// 生成认证令牌
func (s *AuthService) generateToken(user *data.User) (*LoginResponse, error) {
	secretKey := config.GetString("auth.secret_key")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
	}

	token := base64.StdEncoding.EncodeToString([]byte(user.Username + ":" + time.Now().Format(time.RFC3339)))

	return &LoginResponse{
		Token:     token,
		UserID:    user.ID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

// 添加通行密钥
func (s *AuthService) AddPasskey(userID int, credentialID, publicKey string) error {
	credential := &PasskeyCredential{
		ID:        credentialID,
		PublicKey: publicKey,
		UserID:    userID,
	}

	_, err := s.engine.Insert(credential)
	return err
}

// 删除通行密钥
func (s *AuthService) RemovePasskey(userID, credentialID string) error {
	_, err := s.engine.Where("user_id = ? AND id = ?", userID, credentialID).Delete(&PasskeyCredential{})
	return err
}

// 获取用户所有通行密钥
func (s *AuthService) GetPasskeys(userID int) ([]PasskeyCredential, error) {
	var credentials []PasskeyCredential
	err := s.engine.Where("user_id = ?", userID).Find(&credentials)
	return credentials, err
}

// 生成通行密钥认证挑战
func GeneratePasskeyChallenge() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Harbor"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return base64.StdEncoding.EncodeToString(certPEM), base64.StdEncoding.EncodeToString(keyPEM), nil
}

// 验证通行密钥签名
func VerifyPasskeySignature(challenge, signature, publicKeyPEM string) (bool, error) {
	pubKeyBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pubKeyBlock == nil {
		return false, errors.New("解码公钥失败")
	}

	pub, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		return false, err
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}

	chalBytes, err := base64.StdEncoding.DecodeString(challenge)
	if err != nil {
		return false, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("无效的密钥类型")
	}

	err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, chalBytes, sigBytes)
	return err == nil, nil
}

// 认证中间件处理器
type AuthHandler struct{}

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
			c.JSON(401, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// 令牌载荷结构体
type TokenPayload struct {
	UserID   int    `json:"user_id"`  // 用户ID
	Username string `json:"username"` // 用户名
	IsAdmin  bool   `json:"is_admin"` // 是否管理员
}

// 解析认证令牌
func ParseToken(token string) (*TokenPayload, error) {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var payload TokenPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

// 生成认证令牌
func GenerateToken(payload *TokenPayload, secretKey string) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
