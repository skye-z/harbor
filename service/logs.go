package service

import (
	"strconv"

	"github.com/skye-z/harbor/model"
	"github.com/skye-z/harbor/util"

	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

type LogsService struct {
	LogModel model.LogModel
}

func NewLogsService(engine *xorm.Engine) *LogsService {
	ls := new(LogsService)
	ls.LogModel = model.LogModel{
		DB: engine,
	}
	return ls
}

func (ls LogsService) GetNumber(ctx *gin.Context) {
	number, err := ls.LogModel.GetLogNumber()
	if err != nil {
		util.ReturnMessage(ctx, false, "获取日志数量失败")
	} else {
		util.ReturnData(ctx, true, number)
	}
}

func (ls LogsService) GetLogs(ctx *gin.Context) {
	page := ctx.Query("page")
	iPage, err1 := strconv.Atoi(page)
	num := ctx.Query("number")
	iNum, err2 := strconv.Atoi(num)
	if err1 != nil || err2 != nil {
		util.ReturnError(ctx, util.Errors.ParamIllegalError)
		return
	}
	if iNum == 0 {
		iNum = 10
	}
	list, err := ls.LogModel.GetLogs(iPage, iNum)
	if err != nil {
		util.ReturnMessage(ctx, false, "获取日志列表失败")
	} else {
		util.ReturnData(ctx, true, list)
	}
}
