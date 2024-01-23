package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"harbor/model"
	"harbor/util"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"xorm.io/xorm"
)

const IssuerName = "Skye>Quest.Auth"
const tokenKey = "token.secret"

func AuthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		code := ctx.Request.Header.Get("Authorization")
		if code == "" {
			util.ReturnError(ctx, util.Errors.NotLoginError)
			return
		} else if strings.Contains(code, " ") {
			code = code[strings.Index(code, " ")+1:]
		}

		info := jwt.MapClaims{}
		// 密钥
		secret := util.GetString(tokenKey)
		token, err := jwt.ParseWithClaims(code, &info, func(token *jwt.Token) (interface{}, error) {
			key, err := base64.StdEncoding.DecodeString(secret)
			return key, err
		})
		if err != nil {
			util.ReturnError(ctx, util.Errors.TokenNotAvailableError)
			return
		}
		if !token.Valid {
			util.ReturnError(ctx, util.Errors.TokenInvalidError)
			return
		}
		iss, err := info.GetIssuer()
		if err != nil {
			util.ReturnError(ctx, util.Errors.TokenNotAvailableError)
			return
		}
		if iss != IssuerName {
			util.ReturnError(ctx, util.Errors.TokenIllegalError)
			return
		}
		sub, err := info.GetSubject()
		if err != nil {
			util.ReturnError(ctx, util.Errors.TokenNotAvailableError)
			return
		}
		subs := strings.Split(sub, "@")
		uid, err := strconv.ParseInt(subs[1], 10, 64)
		if err != nil {
			util.ReturnError(ctx, util.Errors.TokenNotAvailableError)
			return
		}
		user := model.User{
			Id:    uid,
			Name:  subs[0],
			Admin: subs[2] == "true",
		}
		ctx.Set("user", user)
	}
}

func GenerateToken(user *model.User) (string, int64, error) {
	// 密钥
	secret := util.GetString(tokenKey)
	// 有效小时
	expTime := util.GetInt("token.exp")
	// 过期时间
	exp := time.Now().Add(time.Hour * time.Duration(expTime)).Unix()
	tc := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": exp,
			"iss": IssuerName,
			"sub": fmt.Sprintf("%s@%v@%v", user.Name, user.Id, user.Admin),
		},
	)
	key, _ := base64.StdEncoding.DecodeString(secret)
	token, err := tc.SignedString(key)
	return token, exp, err
}

func ValidateToken(code string) (bool, int, string, error) {
	info := jwt.MapClaims{}
	// 密钥
	secret := util.GetString(tokenKey)
	token, err := jwt.ParseWithClaims(code, &info, func(token *jwt.Token) (interface{}, error) {
		key, err := base64.StdEncoding.DecodeString(secret)
		return key, err
	})
	if err != nil {
		return false, 0, "", err
	}
	if !token.Valid {
		return false, 0, "", nil
	}
	sub, err := info.GetSubject()
	if err != nil {
		return false, 0, "", err
	}
	subs := strings.Split(sub, "@")
	uid, err := strconv.Atoi(subs[1])
	if err != nil {
		return false, 0, "", err
	}
	return true, uid, subs[0], nil
}

type AuthService struct {
	Config *oauth2.Config
	DB     *xorm.Engine
}

func NewAuthService(db *xorm.Engine) *AuthService {
	as := new(AuthService)
	as.DB = db
	as.Config = GetOAuth2Config()
	if as.Config == nil {
		return nil
	}
	return as
}

func GetOAuth2Config() *oauth2.Config {
	if !util.GetBool("oauth2.enable") {
		return nil
	}
	scopes := util.GetString("oauth2.scopes")
	return &oauth2.Config{
		ClientID:     util.GetString("oauth2.clientId"),
		ClientSecret: util.GetString("oauth2.clientSecret"),
		RedirectURL:  util.GetString("oauth2.redirectUrl"),
		Scopes:       strings.Split(scopes, ","),
		Endpoint: oauth2.Endpoint{
			AuthURL:  util.GetString("oauth2.authUrl"),
			TokenURL: util.GetString("oauth2.tokenUrl"),
		},
	}
}

func (as AuthService) Login(c *gin.Context) {
	// 重定向到提供商的授权页面
	url := as.Config.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (as AuthService) Callback(ctx *gin.Context) {
	// 处理提供商的回调并获取访问令牌
	code := ctx.Query("code")
	res, err := as.Config.Exchange(ctx, code)
	if err != nil {
		// 授权服务不可用
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/#/auth/error?code=1")
		return
	}
	// 换取授权信息
	token := res.AccessToken
	id, _, err := as.QueryUserInfo(token)
	if err != nil || id == "" {
		// 授权信息无效
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/#/auth/error?code=2")
		return
	}
	// 获取用户信息
	um := &model.UserModel{
		DB: as.DB,
	}
	user, err := um.GetOAuthUser(id)
	if err != nil {
		// 查询绑定用户失败
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/#/auth/error?code=3")
		return
	}
	if user == nil {
		// 授权用户不存在
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/#/auth/error?code=4")
		return
	}
	// 签发令牌
	token, exp, err := GenerateToken(user)
	if err != nil {
		// 令牌签发失败
		ctx.Redirect(http.StatusTemporaryRedirect, "/app/#/auth/error?code=5")
		return
	}
	logger := &model.LogModel{
		DB: as.DB,
	}
	logger.AddLog("login", "oauth2", ctx.ClientIP())
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/app/#/auth/jump?code=%s&exp=%v", token, exp))
}

func (as AuthService) QueryUserInfo(token string) (string, string, error) {
	// 创建 HTTP 请求
	req, err := http.NewRequest("GET", util.GetString("oauth2.userurl"), nil)
	if err != nil {
		return "", "", err
	}

	// 设置 Authorization 头部，带上 Bearer Token
	req.Header.Set("Authorization", "Bearer "+token)

	// 发起 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// 使用 map 解析 JSON 数据
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", err
	}

	idKey := util.GetString("oauth2.userIdKey")
	nameKey := util.GetString("oauth2.userNameKey")

	oauth2Id := ""
	oauth2Name := ""

	for key, value := range result {
		if key == idKey {
			if strVal, ok := value.(string); ok {
				oauth2Id = strVal
			} else if strVal, ok := value.(int16); ok {
				oauth2Id = strconv.FormatInt(int64(strVal), 10)
			} else if strVal, ok := value.(float64); ok {
				oauth2Id = strconv.FormatFloat(float64(strVal), 'f', -1, 64)
			}
		}
		if key == nameKey {
			if strVal, ok := value.(string); ok {
				oauth2Name = strVal
			}
		}
	}
	return oauth2Id, oauth2Name, err
}
