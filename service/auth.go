package service

import (
	"encoding/base64"
	"fmt"
	"harbor/model"
	"harbor/util"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			"sub": fmt.Sprintf("%s@%v@%v@%v", user.Name, user.Id, user.Admin),
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
