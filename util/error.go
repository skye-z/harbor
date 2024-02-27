/*
公共错误工具

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import "fmt"

type CustomError struct {
	State   bool   `json:"state"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("message: %s, code: %d", e.Message, e.Code)
}

type CustomErrors struct {
	// 未登录
	NotLoginError CustomError
	// 令牌无效
	TokenInvalidError CustomError
	// 令牌不合法
	TokenIllegalError CustomError
	// 令牌不合法
	TokenNotAvailableError CustomError
	// 权限不足
	PermissionDeniedError CustomError
	// 参数为空
	ParamEmptyError CustomError
	// 参数不合法
	ParamIllegalError CustomError
	// 意料之外
	UnexpectedError CustomError
}

var Errors = CustomErrors{
	NotLoginError:          CustomError{false, 10100, "请先登录"},
	TokenInvalidError:      CustomError{false, 10101, "访问令牌无效"},
	TokenIllegalError:      CustomError{false, 10102, "访问令牌不合法"},
	TokenNotAvailableError: CustomError{false, 10103, "令牌无法解析"},
	PermissionDeniedError:  CustomError{false, 10104, "权限不足"},
	ParamEmptyError:        CustomError{false, 10105, "缺少关键参数"},
	ParamIllegalError:      CustomError{false, 10106, "参数类型错误"},
	UnexpectedError:        CustomError{false, 99999, "发生意料之外的错误"},
}
